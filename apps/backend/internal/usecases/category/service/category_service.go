package service

import (
	"context"
	"errors"
	"time"

	"go_refine_dashboard_be/internal/domain/category/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/infrastructure/cache"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/category/dto"
)

type CategoryService interface {
	GetCategories(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string, parentID *uint) ([]dto.CategoryResponse, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error
	BulkUpdateStatus(ctx context.Context, req *dto.BulkUpdateCategoryStatusRequest) error
	ReorderCategories(ctx context.Context, req *dto.ReorderCategoriesRequest) error
}

type categoryService struct {
	categoryRepo    repository.CategoryRepository
	publicPageCache cache.PublicPageCacheInvalidator
}

func NewCategoryService(categoryRepo repository.CategoryRepository, publicPageCache cache.PublicPageCacheInvalidator) CategoryService {
	return &categoryService{
		categoryRepo:    categoryRepo,
		publicPageCache: publicPageCache,
	}
}

func (s *categoryService) GetCategories(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string, parentID *uint) ([]dto.CategoryResponse, int64, error) {
	categories, total, err := s.categoryRepo.FindAll(ctx, skip, limit, sortField, sortOrder, query, status, parentID)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.CategoryResponse
	for _, cat := range categories {
		count, _ := s.categoryRepo.CountProductsByCategoryID(ctx, cat.ID)
		responses = append(responses, mapEntityToResponse(&cat, count))
	}

	return responses, total, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (*dto.CategoryResponse, error) {
	cat, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	count, _ := s.categoryRepo.CountProductsByCategoryID(ctx, cat.ID)
	res := mapEntityToResponse(cat, count)
	return &res, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	exists, err := s.categoryRepo.CheckSlugExists(ctx, req.Slug, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("slug đã tồn tại")
	}

	cat := &entity.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}
	if cat.Status == "" {
		cat.Status = "ACTIVE"
	}

	if err := s.categoryRepo.Create(ctx, cat); err != nil {
		return nil, err
	}

	res := mapEntityToResponse(cat, 0)
	return &res, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	cat, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Slug != nil && *req.Slug != cat.Slug {
		exists, err := s.categoryRepo.CheckSlugExists(ctx, *req.Slug, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("slug đã tồn tại")
		}
		cat.Slug = *req.Slug
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, errors.New("danh mục cha không được là chính nó")
		}
		// Check circular dependency
		currParentID := req.ParentID
		for currParentID != nil {
			if *currParentID == id {
				return nil, errors.New("phát hiện vòng lặp danh mục (circular dependency)")
			}
			parentCat, err := s.categoryRepo.FindByID(ctx, *currParentID)
			if err != nil || parentCat == nil {
				break
			}
			currParentID = parentCat.ParentID
		}
		cat.ParentID = req.ParentID
	}

	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Description != nil {
		cat.Description = *req.Description
	}
	if req.ImageURL != nil {
		cat.ImageURL = *req.ImageURL
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		cat.Status = *req.Status
	}

	if err := s.categoryRepo.Update(ctx, cat); err != nil {
		return nil, err
	}
	s.publicPageCache.InvalidateSource(ctx, pageEntity.PageSectionItemTypeCategory, cat.ID)

	count, _ := s.categoryRepo.CountProductsByCategoryID(ctx, cat.ID)
	res := mapEntityToResponse(cat, count)
	return &res, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	count, err := s.categoryRepo.CountProductsByCategoryID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("không thể xóa danh mục đang có sản phẩm liên kết")
	}
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}
	s.publicPageCache.InvalidateSource(ctx, pageEntity.PageSectionItemTypeCategory, id)
	return nil
}

func (s *categoryService) BulkUpdateStatus(ctx context.Context, req *dto.BulkUpdateCategoryStatusRequest) error {
	if err := s.categoryRepo.BulkUpdateStatus(ctx, req.IDs, req.Status); err != nil {
		return err
	}
	s.publicPageCache.InvalidateSource(ctx, pageEntity.PageSectionItemTypeCategory, req.IDs...)
	return nil
}

func (s *categoryService) ReorderCategories(ctx context.Context, req *dto.ReorderCategoriesRequest) error {
	var items []struct {
		ID        uint
		SortOrder int
	}
	for _, reqItem := range req.Items {
		items = append(items, struct {
			ID        uint
			SortOrder int
		}{ID: reqItem.ID, SortOrder: reqItem.SortOrder})
	}
	return s.categoryRepo.Reorder(ctx, items)
}

func mapEntityToResponse(cat *entity.Category, productCount int64) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		ParentID:     cat.ParentID,
		Description:  cat.Description,
		ImageURL:     cat.ImageURL,
		SortOrder:    cat.SortOrder,
		Status:       cat.Status,
		ProductCount: productCount,
		CreatedAt:    cat.CreatedAt.Format(time.RFC3339),
	}
}
