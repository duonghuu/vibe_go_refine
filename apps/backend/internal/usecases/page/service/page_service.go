package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/infrastructure/cache"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/page/dto"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrPageInvalidInput = errors.New("dữ liệu trang không hợp lệ")
	ErrPageSlugConflict = errors.New("slug trang đã tồn tại")
	ErrPageNotFound     = errors.New("không tìm thấy trang")
)

var pageSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type PageService interface {
	CreatePage(ctx context.Context, authorID uint, req *dto.CreatePageRequest) (*dto.PageResponse, error)
	GetPages(ctx context.Context, query dto.GetPagesQuery) (*dto.PaginatedPageResponse, error)
	GetPageByID(ctx context.Context, id uint) (*dto.PageResponse, error)
	UpdatePage(ctx context.Context, id uint, req *dto.UpdatePageRequest) (*dto.PageResponse, error)
	DeletePage(ctx context.Context, id uint) error
}

type pageService struct {
	pageRepo        repository.PageRepository
	redisClient     *redis.Client
	publicPageCache cache.PublicPageCacheInvalidator
}

func NewPageService(pageRepo repository.PageRepository, redisClient *redis.Client, publicPageCache cache.PublicPageCacheInvalidator) PageService {
	return &pageService{pageRepo: pageRepo, redisClient: redisClient, publicPageCache: publicPageCache}
}

func (s *pageService) invalidateCache(ctx context.Context, patterns ...string) {
	if s.redisClient == nil {
		return
	}

	for _, pattern := range patterns {
		iterator := s.redisClient.Scan(ctx, 0, pattern, 100).Iterator()
		for iterator.Next(ctx) {
			if err := s.redisClient.Del(ctx, iterator.Val()).Err(); err != nil {
				log.Printf("page cache invalidation failed for key %q: %v", iterator.Val(), err)
			}
		}
		if err := iterator.Err(); err != nil {
			log.Printf("page cache scan failed for pattern %q: %v", pattern, err)
		}
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

	s.invalidateCache(ctx, "admin:pages:list:*")
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

func (s *pageService) GetPages(ctx context.Context, query dto.GetPagesQuery) (*dto.PaginatedPageResponse, error) {
	current := query.Current
	if current < 1 {
		current = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	status := strings.TrimSpace(query.Status)
	if status != "" && status != string(entity.PageStatusDraft) && status != string(entity.PageStatusPublished) {
		return nil, fmt.Errorf("%w: trạng thái không hợp lệ", ErrPageInvalidInput)
	}

	sortColumn := "updated_at"
	switch query.SortBy {
	case "id", "title", "slug", "updated_at":
		sortColumn = query.SortBy
	}
	sortOrder := "DESC"
	if strings.EqualFold(query.Order, "asc") {
		sortOrder = "ASC"
	}

	pages, total, err := s.pageRepo.FindAndCount(ctx, repository.PageListFilter{
		TitleLike: strings.TrimSpace(query.Search),
		Status:    entity.PageStatus(status),
		Offset:    (current - 1) * pageSize,
		Limit:     pageSize,
		Sort:      sortColumn + " " + sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy danh sách trang: %w", err)
	}

	data := make([]dto.PageListItemResponse, 0, len(pages))
	for _, page := range pages {
		data = append(data, dto.PageListItemResponse{
			ID:        page.ID,
			Title:     page.Title,
			Slug:      page.Slug,
			Status:    page.Status,
			AuthorID:  page.AuthorID,
			CreatedAt: page.CreatedAt,
			UpdatedAt: page.UpdatedAt,
		})
	}

	return &dto.PaginatedPageResponse{Data: data, Total: total}, nil
}

func (s *pageService) GetPageByID(ctx context.Context, id uint) (*dto.PageResponse, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPageNotFound
		}
		return nil, fmt.Errorf("lỗi lấy chi tiết trang: %w", err)
	}

	return mapPageToResponse(page), nil
}

func (s *pageService) UpdatePage(ctx context.Context, id uint, req *dto.UpdatePageRequest) (*dto.PageResponse, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPageNotFound
		}
		return nil, fmt.Errorf("lỗi lấy trang cần cập nhật: %w", err)
	}

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

	exists, err := s.pageRepo.CheckSlugExistsExceptID(ctx, slug, id)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra slug trang: %w", err)
	}
	if exists {
		return nil, ErrPageSlugConflict
	}

	oldSlug := page.Slug
	page.Title = title
	page.Slug = slug
	page.Content = content
	page.Status = entity.PageStatus(status)
	if err := s.pageRepo.Update(ctx, page); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") && strings.Contains(strings.ToLower(err.Error()), "slug") {
			return nil, ErrPageSlugConflict
		}
		return nil, fmt.Errorf("không thể cập nhật trang: %w", err)
	}

	s.invalidateCache(ctx,
		"admin:pages:list:*",
		fmt.Sprintf("seo:page:%d", id),
	)
	s.publicPageCache.InvalidateSlugs(ctx, oldSlug, slug)
	return mapPageToResponse(page), nil
}

func (s *pageService) DeletePage(ctx context.Context, id uint) error {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPageNotFound
		}
		return fmt.Errorf("lỗi tìm trang cần xóa: %w", err)
	}

	if err := s.pageRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("không thể xóa trang: %w", err)
	}

	s.invalidateCache(ctx,
		"admin:pages:list:*",
		fmt.Sprintf("seo:page:%d", id),
	)
	s.publicPageCache.InvalidateSlugs(ctx, page.Slug)
	return nil
}

func mapPageToResponse(page *entity.Page) *dto.PageResponse {
	return &dto.PageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Slug:      page.Slug,
		Content:   page.Content,
		Status:    string(page.Status),
		AuthorID:  page.AuthorID,
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
	}
}
