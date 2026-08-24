package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go_refine_dashboard_be/internal/domain/postcategory/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/postcategory/dto"
)

type PostCategoryService interface {
	GetPostCategoryTree(ctx context.Context, typeCode string) ([]dto.PostCategoryTreeResponse, error)
	CreatePostCategory(ctx context.Context, req *dto.CreatePostCategoryRequest) (*dto.PostCategoryResponse, error)
	GetPostCategoryByID(ctx context.Context, id uint) (*dto.PostCategoryResponse, error)
	UpdatePostCategory(ctx context.Context, id uint, req *dto.UpdatePostCategoryRequest) (*dto.PostCategoryResponse, error)
	GetPostCategories(ctx context.Context, skip, limit int, sortField, sortOrder, query, status, typeCode string) (*dto.PaginatedPostCategoryResponse, error)
	DeletePostCategory(ctx context.Context, id uint) error
}

type postCategoryService struct {
	repo        repository.PostCategoryRepository
	redisClient *redis.Client
}

func NewPostCategoryService(repo repository.PostCategoryRepository, redisClient *redis.Client) PostCategoryService {
	return &postCategoryService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *postCategoryService) invalidateCache(ctx context.Context, typeCode string) {
	cacheKey := fmt.Sprintf("techbite:post_categories:tree:%s", typeCode)
	s.redisClient.Del(ctx, cacheKey)
}

func mapEntityToResponse(e *entity.PostCategory) *dto.PostCategoryResponse {
	return &dto.PostCategoryResponse{
		ID:          e.ID,
		TypeCode:    e.TypeCode,
		Name:        e.Name,
		Slug:        e.Slug,
		ParentID:    e.ParentID,
		Description: e.Description,
		ImageURL:    e.ImageURL,
		SortOrder:   e.SortOrder,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
	}
}

func (s *postCategoryService) GetPostCategoryTree(ctx context.Context, typeCode string) ([]dto.PostCategoryTreeResponse, error) {
	items, err := s.repo.FindTreeByTypeCode(ctx, typeCode)
	if err != nil {
		return nil, err
	}

	var res []dto.PostCategoryTreeResponse
	for _, item := range items {
		res = append(res, dto.PostCategoryTreeResponse{
			ID:       item.ID,
			Name:     item.Name,
			TypeCode: item.TypeCode,
		})
	}
	if res == nil {
		res = []dto.PostCategoryTreeResponse{}
	}
	return res, nil
}

func (s *postCategoryService) CreatePostCategory(ctx context.Context, req *dto.CreatePostCategoryRequest) (*dto.PostCategoryResponse, error) {
	existing, err := s.repo.FindBySlug(ctx, req.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("Slug đã tồn tại")
	}

	if req.ParentID != nil {
		parentCategory, err := s.repo.FindByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Danh mục cha không tồn tại")
			}
			return nil, err
		}
		if parentCategory.TypeCode != req.TypeCode {
			return nil, errors.New("Danh mục cha không thuộc cùng loại bài viết")
		}
	}

	status := "ACTIVE"
	if req.Status != "" {
		status = req.Status
	}

	item := &entity.PostCategory{
		TypeCode:    req.TypeCode,
		Name:        req.Name,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		Status:      status,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, item.TypeCode)
	return mapEntityToResponse(item), nil
}

func (s *postCategoryService) GetPostCategoryByID(ctx context.Context, id uint) (*dto.PostCategoryResponse, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy danh mục bài viết")
		}
		return nil, err
	}
	return mapEntityToResponse(item), nil
}

func (s *postCategoryService) UpdatePostCategory(ctx context.Context, id uint, req *dto.UpdatePostCategoryRequest) (*dto.PostCategoryResponse, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy danh mục bài viết")
		}
		return nil, err
	}

	slugExisting, err := s.repo.FindBySlug(ctx, req.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if slugExisting != nil && slugExisting.ID != id {
		return nil, errors.New("Slug đã tồn tại")
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, errors.New("Không thể chọn danh mục hiện tại làm danh mục cha")
		}
		
		parentCategory, err := s.repo.FindByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Danh mục cha không tồn tại")
			}
			return nil, err
		}
		if parentCategory.TypeCode != existing.TypeCode {
			return nil, errors.New("Danh mục cha không thuộc cùng loại bài viết")
		}
		
		// To prevent circular reference fully, we'd need to check if parent is a descendant of current node.
		// For simplicity, we just rely on Frontend to block it or assume depth 1 check.
	}

	existing.Name = req.Name
	existing.Slug = req.Slug
	existing.ParentID = req.ParentID
	existing.Description = req.Description
	existing.ImageURL = req.ImageURL
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.SortOrder = req.SortOrder

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, existing.TypeCode)
	return mapEntityToResponse(existing), nil
}

func (s *postCategoryService) GetPostCategories(ctx context.Context, skip, limit int, sortField, sortOrder, query, status, typeCode string) (*dto.PaginatedPostCategoryResponse, error) {
	items, total, err := s.repo.FindAll(ctx, skip, limit, sortField, sortOrder, query, status, typeCode)
	if err != nil {
		return nil, err
	}

	res := make([]dto.PostCategoryResponse, 0, len(items))
	for _, item := range items {
		res = append(res, *mapEntityToResponse(&item))
	}

	return &dto.PaginatedPostCategoryResponse{
		Data:  res,
		Total: total,
	}, nil
}

func (s *postCategoryService) DeletePostCategory(ctx context.Context, id uint) error {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.invalidateCache(ctx, category.TypeCode)
	return nil
}
