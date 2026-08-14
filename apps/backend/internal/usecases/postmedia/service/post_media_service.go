package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go_refine_dashboard_be/internal/domain/post/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/postmedia/dto"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PostMediaService interface {
	GetPostMedia(ctx context.Context, postID uint, collection string) ([]dto.PostMediaResponse, int64, error)
	CreatePostMedia(ctx context.Context, postID uint, req *dto.CreatePostMediaRequest) (*dto.PostMediaResponse, error)
	SyncPostMedia(ctx context.Context, postID uint, collection string, req *dto.SyncPostMediaRequest) error
	DeletePostMedia(ctx context.Context, postID uint, mediaID uint, collection string) error
}

type postMediaService struct {
	repo        repository.PostMediaRepository
	redisClient *redis.Client
}

func NewPostMediaService(repo repository.PostMediaRepository, redisClient *redis.Client) PostMediaService {
	return &postMediaService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *postMediaService) getCacheKey(postID uint, collection string) string {
	if collection == "" {
		return fmt.Sprintf("cache:post:%d:media:all", postID)
	}
	return fmt.Sprintf("cache:post:%d:media:%s", postID, collection)
}

func (s *postMediaService) invalidateCache(ctx context.Context, postID uint) {
	// Delete all caches related to this post
	keys, err := s.redisClient.Keys(ctx, fmt.Sprintf("cache:post:%d:media:*", postID)).Result()
	if err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}
}

func (s *postMediaService) GetPostMedia(ctx context.Context, postID uint, collection string) ([]dto.PostMediaResponse, int64, error) {
	cacheKey := s.getCacheKey(postID, collection)
	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var res []dto.PostMediaResponse
		if json.Unmarshal([]byte(cachedData), &res) == nil {
			return res, int64(len(res)), nil
		}
	}

	items, total, err := s.repo.FindByPostID(ctx, postID, collection)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.PostMediaResponse
	for _, item := range items {
		respItem := dto.PostMediaResponse{
			ID:         item.ID,
			PostID:     item.PostID,
			MediaID:    item.MediaID,
			Collection: item.Collection,
			SortOrder:  item.SortOrder,
		}
		// Load media information
		media, mediaErr := s.repo.GetMediaByID(ctx, item.MediaID)
		if mediaErr == nil {
			respItem.Media = media
		}
		res = append(res, respItem)
	}

	if res == nil {
		res = []dto.PostMediaResponse{}
	}

	// Cache the result
	if cacheBytes, err := json.Marshal(res); err == nil {
		s.redisClient.SetEx(ctx, cacheKey, string(cacheBytes), 24*time.Hour)
	}

	return res, total, nil
}

func (s *postMediaService) CreatePostMedia(ctx context.Context, postID uint, req *dto.CreatePostMediaRequest) (*dto.PostMediaResponse, error) {
	// Verify media exists
	media, err := s.repo.GetMediaByID(ctx, req.MediaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Media không tồn tại")
		}
		return nil, err
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	item := &entity.PostMedia{
		PostID:     postID,
		MediaID:    req.MediaID,
		Collection: req.Collection,
		SortOrder:  sortOrder,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, postID)

	return &dto.PostMediaResponse{
		ID:         item.ID,
		PostID:     item.PostID,
		MediaID:    item.MediaID,
		Collection: item.Collection,
		SortOrder:  item.SortOrder,
		Media:      media,
	}, nil
}

func (s *postMediaService) SyncPostMedia(ctx context.Context, postID uint, collection string, req *dto.SyncPostMediaRequest) error {
	var items []entity.PostMedia
	for _, m := range req.Media {
		items = append(items, entity.PostMedia{
			PostID:     postID,
			MediaID:    m.ID,
			Collection: collection,
			SortOrder:  m.SortOrder,
		})
	}

	if err := s.repo.SyncByCollection(ctx, postID, collection, items); err != nil {
		return err
	}

	s.invalidateCache(ctx, postID)
	return nil
}

func (s *postMediaService) DeletePostMedia(ctx context.Context, postID uint, mediaID uint, collection string) error {
	if err := s.repo.Delete(ctx, postID, mediaID, collection); err != nil {
		return err
	}

	s.invalidateCache(ctx, postID)
	return nil
}
