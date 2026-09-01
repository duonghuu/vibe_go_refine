package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go_refine_dashboard_be/internal/domain/media/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/media/dto"
)

var (
	ErrMediaNotFound  = errors.New("media không tồn tại")
	ErrMediaForbidden = errors.New("không có quyền xóa media này")
	ErrMediaAttached  = errors.New("media đang được liên kết, không thể xóa")
)

type MediaService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userId uint) (*dto.UploadMediaResponse, error)
	DeleteFile(ctx context.Context, id uint, userId uint, role string) error
	ListMedia(ctx context.Context, query dto.GetMediaListQuery, userID uint, role string) (*dto.MediaListResponse, error)
}

func (s *mediaService) ListMedia(ctx context.Context, query dto.GetMediaListQuery, userID uint, role string) (*dto.MediaListResponse, error) {
	current := query.Current
	if current < 1 {
		current = 1
	}
	size := query.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	rows, total, err := s.mediaRepo.FindAndCount(ctx, userID, role == "ADMIN", query.Q, query.MimeType, query.Status, (current-1)*size, size)
	if err != nil {
		return nil, err
	}
	data := make([]dto.MediaListItemResponse, 0, len(rows))
	for _, m := range rows {
		data = append(data, dto.MediaListItemResponse{ID: m.ID, FileName: m.FileName, OriginalURL: m.OriginalUrl, ThumbnailURL: m.ThumbnailUrl, MediumURL: m.MediumUrl, MimeType: m.MimeType, Size: m.Size, Status: string(m.Status), CreatedAt: m.CreatedAt})
	}
	return &dto.MediaListResponse{Data: data, Total: total}, nil
}

type mediaService struct {
	mediaRepo repository.MediaRepository
}

func NewMediaService(mediaRepo repository.MediaRepository) MediaService {
	return &mediaService{
		mediaRepo: mediaRepo,
	}
}

func (s *mediaService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, userId uint) (*dto.UploadMediaResponse, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 1. Read first 512 bytes for Magic Bytes detection
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, errors.New("failed to read file bytes")
	}

	// Rewind the file to the beginning
	src.Seek(0, 0)

	// Detect content type
	contentType := http.DetectContentType(buffer[:n])

	// Only allow images (can be customized)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return nil, errors.New("invalid file type: only JPEG, PNG, WEBP are allowed")
	}

	// 2. Generate secure file name
	ext := filepath.Ext(fileHeader.Filename)
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	newFileName := fmt.Sprintf("%s_%d%s", hex.EncodeToString(randomBytes), time.Now().Unix(), ext)

	// 3. Save to disk
	uploadDir := "public/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, errors.New("failed to create upload directory")
	}

	filePath := filepath.Join(uploadDir, newFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.New("failed to save file")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, errors.New("failed to write file")
	}

	// Mock URLs (in production, we might upload to S3 or generate absolute URLs)
	originalUrl := "/uploads/" + newFileName

	// 4. Create database record
	media := &entity.Media{
		OriginalName: fileHeader.Filename,
		FileName:     newFileName,
		MimeType:     contentType,
		Size:         fileHeader.Size,
		OriginalUrl:  originalUrl,
		Status:       entity.MediaStatusTemporary,
		OwnerID:      userId, // Currently hardcoded in controller due to lack of auth
	}

	if err := s.mediaRepo.Create(ctx, media); err != nil {
		// Clean up the file if DB save fails
		os.Remove(filePath)
		return nil, err
	}

	return &dto.UploadMediaResponse{
		ID:           media.ID,
		FileName:     media.FileName,
		OriginalUrl:  media.OriginalUrl,
		ThumbnailUrl: media.ThumbnailUrl,
		MediumUrl:    media.MediumUrl,
		Status:       string(media.Status),
	}, nil
}

func (s *mediaService) DeleteFile(ctx context.Context, id uint, userId uint, role string) error {
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return ErrMediaNotFound
	}

	if media.Status == entity.MediaStatusAttached {
		return ErrMediaAttached
	}
	if role != "ADMIN" && media.OwnerID != userId {
		return ErrMediaForbidden
	}

	// Delete from DB
	if err := s.mediaRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete from Disk
	uploadDir := "public/uploads"
	filePath := filepath.Join(uploadDir, media.FileName)
	os.Remove(filePath) // Ignore error if file already gone

	return nil
}
