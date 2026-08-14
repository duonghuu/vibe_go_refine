package service

import (
	"context"
	"fmt"

	"go_refine_dashboard_be/internal/domain/post/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/postmeta/dto"

	"github.com/redis/go-redis/v9"
)

type PostMetaService interface {
	GetByPostID(ctx context.Context, req *dto.GetPostMetaRequest) (*dto.GetPostMetaResponse, error)
	Sync(ctx context.Context, req *dto.SyncPostMetaRequest) (*dto.SyncPostMetaResponse, error)
	Delete(ctx context.Context, req *dto.DeletePostMetaRequest) (*dto.DeletePostMetaResponse, error)
}

type postMetaService struct {
	repo  repository.PostMetaRepository
	redis *redis.Client
}

func NewPostMetaService(repo repository.PostMetaRepository, redis *redis.Client) PostMetaService {
	return &postMetaService{
		repo:  repo,
		redis: redis,
	}
}

func (s *postMetaService) invalidateCache(ctx context.Context, postID uint) {
	cacheKey := fmt.Sprintf("post:%d:meta", postID)
	_ = s.redis.Del(ctx, cacheKey).Err() // Error is ignored because cache invalidation failure shouldn't fail the request
}

func (s *postMetaService) GetByPostID(ctx context.Context, req *dto.GetPostMetaRequest) (*dto.GetPostMetaResponse, error) {
	items, err := s.repo.FindByPostID(ctx, req.PostID)
	if err != nil {
		return nil, err
	}

	var res []dto.PostMetaResponse
	for _, item := range items {
		res = append(res, dto.PostMetaResponse{
			ID:        item.ID,
			PostID:    item.PostID,
			Key:       item.Key,
			Value:     item.Value,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	if res == nil {
		res = make([]dto.PostMetaResponse, 0)
	}

	return &dto.GetPostMetaResponse{
		Data: res,
	}, nil
}

func (s *postMetaService) Sync(ctx context.Context, req *dto.SyncPostMetaRequest) (*dto.SyncPostMetaResponse, error) {
	var items []entity.PostMeta
	for _, m := range req.Meta {
		items = append(items, entity.PostMeta{
			PostID: req.PostID,
			Key:    m.Key,
			Value:  m.Value,
		})
	}

	if err := s.repo.Sync(ctx, req.PostID, items); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, req.PostID)

	return &dto.SyncPostMetaResponse{
		Message: "Đồng bộ meta thành công",
	}, nil
}

func (s *postMetaService) Delete(ctx context.Context, req *dto.DeletePostMetaRequest) (*dto.DeletePostMetaResponse, error) {
	if err := s.repo.Delete(ctx, req.PostID, req.Key); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, req.PostID)

	return &dto.DeletePostMetaResponse{
		Message: "Xóa meta thành công",
	}, nil
}
