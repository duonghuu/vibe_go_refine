package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/page/dto"

	"github.com/redis/go-redis/v9"
)

var (
	ErrPageInvalidInput = errors.New("dữ liệu trang không hợp lệ")
	ErrPageSlugConflict = errors.New("slug trang đã tồn tại")
)

var pageSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type PageService interface {
	CreatePage(ctx context.Context, authorID uint, req *dto.CreatePageRequest) (*dto.PageResponse, error)
}

type pageService struct {
	pageRepo    repository.PageRepository
	redisClient *redis.Client
}

func NewPageService(pageRepo repository.PageRepository, redisClient *redis.Client) PageService {
	return &pageService{pageRepo: pageRepo, redisClient: redisClient}
}

func (s *pageService) invalidateCache(ctx context.Context) {
	if s.redisClient == nil {
		return
	}

	iterator := s.redisClient.Scan(ctx, 0, "admin:pages:list:*", 100).Iterator()
	for iterator.Next(ctx) {
		if err := s.redisClient.Del(ctx, iterator.Val()).Err(); err != nil {
			log.Printf("page cache invalidation failed for key %q: %v", iterator.Val(), err)
		}
	}
	if err := iterator.Err(); err != nil {
		log.Printf("page cache scan failed: %v", err)
	}
}

func (s *pageService) CreatePage(ctx context.Context, authorID uint, req *dto.CreatePageRequest) (*dto.PageResponse, error) {
	title := strings.TrimSpace(req.Title)
	slug := strings.TrimSpace(req.Slug)
	content := strings.TrimSpace(req.Content)
	status := strings.TrimSpace(req.Status)

	if title == "" || content == "" || !pageSlugPattern.MatchString(slug) {
		return nil, fmt.Errorf("%w: tiêu đề, slug hoặc nội dung không hợp lệ", ErrPageInvalidInput)
	}
	if status != string(entity.PageStatusDraft) && status != string(entity.PageStatusPublished) {
		return nil, fmt.Errorf("%w: trạng thái không hợp lệ", ErrPageInvalidInput)
	}

	exists, err := s.pageRepo.CheckSlugExists(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra slug trang: %w", err)
	}
	if exists {
		return nil, ErrPageSlugConflict
	}

	page := &entity.Page{
		Title:    title,
		Slug:     slug,
		Content:  content,
		Status:   entity.PageStatus(status),
		AuthorID: authorID,
	}
	if err := s.pageRepo.Create(ctx, page); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") && strings.Contains(strings.ToLower(err.Error()), "slug") {
			return nil, ErrPageSlugConflict
		}
		return nil, fmt.Errorf("không thể tạo trang: %w", err)
	}

	s.invalidateCache(ctx)
	return &dto.PageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Slug:      page.Slug,
		Content:   page.Content,
		Status:    string(page.Status),
		AuthorID:  page.AuthorID,
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
	}, nil
}
