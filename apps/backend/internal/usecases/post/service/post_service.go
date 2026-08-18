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
	GetPosts(ctx context.Context, query dto.GetPostsQuery) (*dto.PaginatedPostResponse, error)
	DeletePost(ctx context.Context, id uint) error
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

func (s *postService) GetPosts(ctx context.Context, query dto.GetPostsQuery) (*dto.PaginatedPostResponse, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	sort := ""
	if query.SortBy != "" {
		order := "asc"
		if query.Order != "" {
			order = query.Order
		}
		sort = fmt.Sprintf("%s %s", query.SortBy, order)
	}

	posts, total, err := s.postRepo.FindAndCount(ctx, query.Title, query.TypeCode, offset, pageSize, sort)
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy danh sách bài viết: %v", err)
	}

	var data []dto.PostResponse
	if len(posts) == 0 {
		data = []dto.PostResponse{}
	}
	for _, p := range posts {
		data = append(data, dto.PostResponse{
			ID:        p.ID,
			TypeCode:  p.TypeCode,
			Title:     p.Title,
			Slug:      p.Slug,
			// Content is omitted in list
			AuthorID:  p.AuthorID,
			CreatedAt: p.CreatedAt,
		})
	}

	return &dto.PaginatedPostResponse{
		Data:  data,
		Total: total,
	}, nil
}

func (s *postService) DeletePost(ctx context.Context, id uint) error {
	if err := s.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("không thể xóa bài viết: %v", err)
	}

	// Invalidate Cache for posts list
	pattern := "posts:list:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
	}

	// Note: ideally we would invalidate specific type cache if we query the post first, but this is a broad invalidation
	return nil
}
