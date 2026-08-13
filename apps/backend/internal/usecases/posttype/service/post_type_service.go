package service

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go_refine_dashboard_be/internal/domain/posttype/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/posttype/dto"
)

type PostTypeService interface {
	GetPostTypes(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string) ([]dto.PostTypeResponse, int64, error)
	GetPostTypeByID(ctx context.Context, id uint) (*dto.PostTypeResponse, error)
	CreatePostType(ctx context.Context, req *dto.CreatePostTypeRequest) (*dto.PostTypeResponse, error)
	UpdatePostType(ctx context.Context, id uint, req *dto.UpdatePostTypeRequest) (*dto.PostTypeResponse, error)
	DeletePostType(ctx context.Context, id uint) error
}

type postTypeService struct {
	repo        repository.PostTypeRepository
	redisClient *redis.Client
}

func NewPostTypeService(repo repository.PostTypeRepository, redisClient *redis.Client) PostTypeService {
	return &postTypeService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *postTypeService) invalidateCache(ctx context.Context) {
	s.redisClient.Del(ctx, "techbite:post_types:active_list")
}

func mapEntityToResponse(e *entity.PostType) *dto.PostTypeResponse {
	return &dto.PostTypeResponse{
		ID:        e.ID,
		Code:      e.Code,
		Name:      e.Name,
		Status:    e.Status,
		SortOrder: e.SortOrder,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (s *postTypeService) GetPostTypes(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string) ([]dto.PostTypeResponse, int64, error) {
	items, total, err := s.repo.FindAll(ctx, skip, limit, sortField, sortOrder, query, status)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.PostTypeResponse
	for _, item := range items {
		res = append(res, *mapEntityToResponse(&item))
	}
	if res == nil {
		res = []dto.PostTypeResponse{}
	}
	return res, total, nil
}

func (s *postTypeService) GetPostTypeByID(ctx context.Context, id uint) (*dto.PostTypeResponse, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy loại bài viết")
		}
		return nil, err
	}
	return mapEntityToResponse(item), nil
}

func (s *postTypeService) CreatePostType(ctx context.Context, req *dto.CreatePostTypeRequest) (*dto.PostTypeResponse, error) {
	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("Mã loại bài viết đã tồn tại")
	}

	status := "ACTIVE"
	if req.Status != "" {
		status = req.Status
	}

	item := &entity.PostType{
		Code:      req.Code,
		Name:      req.Name,
		Status:    status,
		SortOrder: req.SortOrder,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx)
	return mapEntityToResponse(item), nil
}

func (s *postTypeService) UpdatePostType(ctx context.Context, id uint, req *dto.UpdatePostTypeRequest) (*dto.PostTypeResponse, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Không tìm thấy loại bài viết")
		}
		return nil, err
	}

	if req.Status != "" {
		item.Status = req.Status
	}
	item.Name = req.Name
	item.SortOrder = req.SortOrder

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx)
	return mapEntityToResponse(item), nil
}

func (s *postTypeService) DeletePostType(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Không tìm thấy loại bài viết")
		}
		return err
	}

	postCount, err := s.repo.CountPostsByPostTypeID(ctx, id)
	if err == nil && postCount > 0 {
		return errors.New("Không thể xóa loại bài viết đang được sử dụng")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.invalidateCache(ctx)
	return nil
}
