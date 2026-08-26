package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go_refine_dashboard_be/internal/domain/media/entity"
	postEntity "go_refine_dashboard_be/internal/domain/post/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/postmedia/dto"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound      = errors.New("bài viết không tồn tại")
	ErrMediaNotFound     = errors.New("media không tồn tại")
	ErrMediaForbidden    = errors.New("không có quyền sử dụng media này")
	ErrCollectionInvalid = errors.New("collection không hợp lệ")
	ErrThumbnailLimit    = errors.New("mỗi bài viết chỉ được có một thumbnail")
	ErrMediaDuplicate    = errors.New("media bị trùng trong collection")
	ErrMediaTypeInvalid  = errors.New("media phải là file hình ảnh")
)

type PostMediaService interface {
	GetPostMedia(ctx context.Context, postID uint, collection string) ([]dto.PostMediaResponse, int64, error)
	CreatePostMedia(ctx context.Context, postID uint, req *dto.CreatePostMediaRequest, userID uint, role string) (*dto.PostMediaResponse, error)
	SyncPostMedia(ctx context.Context, postID uint, collection string, req *dto.SyncPostMediaRequest, userID uint, role string) ([]dto.PostMediaResponse, error)
	DeletePostMedia(ctx context.Context, postID uint, mediaID uint, collection string, userID uint, role string) error
}

type postMediaService struct {
	repo        repository.PostMediaRepository
	redisClient *redis.Client
}

func NewPostMediaService(repo repository.PostMediaRepository, redisClient *redis.Client) PostMediaService {
	return &postMediaService{repo: repo, redisClient: redisClient}
}

func (s *postMediaService) getCacheKey(postID uint, collection string) string {
	if collection == "" {
		return fmt.Sprintf("cache:post:%d:media:all", postID)
	}
	return fmt.Sprintf("cache:post:%d:media:%s", postID, collection)
}

func (s *postMediaService) invalidateCache(ctx context.Context, postID uint) {
	if s.redisClient == nil {
		return
	}
	iterator := s.redisClient.Scan(ctx, 0, fmt.Sprintf("cache:post:%d:media:*", postID), 100).Iterator()
	for iterator.Next(ctx) {
		_ = s.redisClient.Del(ctx, iterator.Val()).Err()
	}
}

func (s *postMediaService) ensurePost(ctx context.Context, postID uint) error {
	exists, err := s.repo.PostExists(ctx, postID)
	if err != nil {
		return fmt.Errorf("kiểm tra bài viết: %w", err)
	}
	if !exists {
		return ErrPostNotFound
	}
	return nil
}

func validateCollection(collection string) error {
	if collection != postEntity.PostMediaCollectionThumbnail && collection != postEntity.PostMediaCollectionGallery {
		return ErrCollectionInvalid
	}
	return nil
}

func canManageMedia(media entity.Media, userID uint, role string) bool {
	return role == "ADMIN" || media.OwnerID == userID
}

func (s *postMediaService) validateMedia(ctx context.Context, mediaIDs []uint, userID uint, role string) ([]entity.Media, error) {
	seen := make(map[uint]struct{}, len(mediaIDs))
	for _, id := range mediaIDs {
		if _, exists := seen[id]; exists {
			return nil, ErrMediaDuplicate
		}
		seen[id] = struct{}{}
	}

	media, err := s.repo.GetMediaByIDs(ctx, mediaIDs)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách media: %w", err)
	}
	if len(media) != len(mediaIDs) {
		return nil, ErrMediaNotFound
	}
	for _, item := range media {
		if !canManageMedia(item, userID, role) {
			return nil, ErrMediaForbidden
		}
		if !strings.HasPrefix(strings.ToLower(item.MimeType), "image/") {
			return nil, ErrMediaTypeInvalid
		}
	}
	return media, nil
}

func mapPostMedia(item postEntity.PostMedia) dto.PostMediaResponse {
	return dto.PostMediaResponse{
		ID:         item.ID,
		PostID:     item.PostID,
		MediaID:    item.MediaID,
		Collection: item.Collection,
		SortOrder:  item.SortOrder,
		Media:      dto.MapMediaResponse(item.Media),
	}
}

func (s *postMediaService) GetPostMedia(ctx context.Context, postID uint, collection string) ([]dto.PostMediaResponse, int64, error) {
	if err := s.ensurePost(ctx, postID); err != nil {
		return nil, 0, err
	}
	if collection != "" {
		if err := validateCollection(collection); err != nil {
			return nil, 0, err
		}
	}

	cacheKey := s.getCacheKey(postID, collection)
	if s.redisClient != nil {
		if cached, err := s.redisClient.Get(ctx, cacheKey).Result(); err == nil {
			var result []dto.PostMediaResponse
			if json.Unmarshal([]byte(cached), &result) == nil {
				return result, int64(len(result)), nil
			}
		}
	}

	items, total, err := s.repo.FindByPostID(ctx, postID, collection)
	if err != nil {
		return nil, 0, fmt.Errorf("lấy media bài viết: %w", err)
	}
	result := make([]dto.PostMediaResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapPostMedia(item))
	}
	if s.redisClient != nil {
		if payload, marshalErr := json.Marshal(result); marshalErr == nil {
			_ = s.redisClient.SetEx(ctx, cacheKey, payload, 24*time.Hour).Err()
		}
	}
	return result, total, nil
}

func (s *postMediaService) CreatePostMedia(ctx context.Context, postID uint, req *dto.CreatePostMediaRequest, userID uint, role string) (*dto.PostMediaResponse, error) {
	if err := s.ensurePost(ctx, postID); err != nil {
		return nil, err
	}
	if err := validateCollection(req.Collection); err != nil {
		return nil, err
	}
	if _, err := s.validateMedia(ctx, []uint{req.MediaID}, userID, role); err != nil {
		return nil, err
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	item := postEntity.PostMedia{PostID: postID, MediaID: req.MediaID, Collection: req.Collection, SortOrder: sortOrder}
	if err := s.repo.CreateAndAttach(ctx, &item); err != nil {
		return nil, fmt.Errorf("liên kết media: %w", err)
	}
	s.invalidateCache(ctx, postID)
	items, _, err := s.GetPostMedia(ctx, postID, req.Collection)
	if err != nil {
		return nil, err
	}
	for _, result := range items {
		if result.MediaID == req.MediaID {
			return &result, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *postMediaService) SyncPostMedia(ctx context.Context, postID uint, collection string, req *dto.SyncPostMediaRequest, userID uint, role string) ([]dto.PostMediaResponse, error) {
	if err := s.ensurePost(ctx, postID); err != nil {
		return nil, err
	}
	if err := validateCollection(collection); err != nil {
		return nil, err
	}
	if collection == postEntity.PostMediaCollectionThumbnail && len(req.Media) > 1 {
		return nil, ErrThumbnailLimit
	}

	ids := make([]uint, 0, len(req.Media))
	items := make([]postEntity.PostMedia, 0, len(req.Media))
	for _, item := range req.Media {
		ids = append(ids, item.ID)
		items = append(items, postEntity.PostMedia{PostID: postID, MediaID: item.ID, Collection: collection, SortOrder: item.SortOrder})
	}
	if _, err := s.validateMedia(ctx, ids, userID, role); err != nil {
		return nil, err
	}
	if err := s.repo.SyncByCollectionAndAttach(ctx, postID, collection, items, ids); err != nil {
		return nil, fmt.Errorf("đồng bộ media: %w", err)
	}
	s.invalidateCache(ctx, postID)
	result, _, err := s.GetPostMedia(ctx, postID, collection)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *postMediaService) DeletePostMedia(ctx context.Context, postID uint, mediaID uint, collection string, userID uint, role string) error {
	if err := s.ensurePost(ctx, postID); err != nil {
		return err
	}
	if collection != "" {
		if err := validateCollection(collection); err != nil {
			return err
		}
	}
	media, err := s.repo.GetMediaByID(ctx, mediaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("lấy media: %w", err)
	}
	if !canManageMedia(*media, userID, role) {
		return ErrMediaForbidden
	}
	if err := s.repo.Delete(ctx, postID, mediaID, collection); err != nil {
		return fmt.Errorf("gỡ liên kết media: %w", err)
	}
	s.invalidateCache(ctx, postID)
	return nil
}
