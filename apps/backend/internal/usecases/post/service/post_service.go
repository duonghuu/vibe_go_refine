package service

import (
	"context"
	"errors"
	"fmt"

	"go_refine_dashboard_be/internal/domain/post/entity"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/post/dto"

	"github.com/redis/go-redis/v9"
)

type PostService interface {
	CreatePost(ctx context.Context, authorID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error)
}

type postService struct {
	postRepo    repository.PostRepository
	redisClient *redis.Client
}

func NewPostService(postRepo repository.PostRepository, redisClient *redis.Client) PostService {
	return &postService{
		postRepo:    postRepo,
		redisClient: redisClient,
	}
}

func (s *postService) CreatePost(ctx context.Context, authorID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	exists, err := s.postRepo.CheckSlugExists(ctx, req.Slug)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra slug: %v", err)
	}
	if exists {
		return nil, errors.New("slug đã tồn tại, vui lòng chọn slug khác")
	}

	post := &entity.Post{
		TypeCode: req.TypeCode,
		Title:    req.Title,
		Slug:     req.Slug,
		Content:  req.Content,
		AuthorID: authorID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("không thể tạo bài viết: %v", err)
	}

	// Invalidate Cache for posts list
	pattern := "posts:list:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}
	
	typePattern := fmt.Sprintf("posts:type:%s:*", req.TypeCode)
	typeKeys, err := s.redisClient.Keys(ctx, typePattern).Result()
	if err == nil && len(typeKeys) > 0 {
		s.redisClient.Del(ctx, typeKeys...)
	}

	return &dto.PostResponse{
		ID:        post.ID,
		TypeCode:  post.TypeCode,
		Title:     post.Title,
		Slug:      post.Slug,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
	}, nil
}
