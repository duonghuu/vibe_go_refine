package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"go_refine_dashboard_be/internal/domain/post/entity"
	"go_refine_dashboard_be/internal/infrastructure/cache"
	"go_refine_dashboard_be/internal/infrastructure/repository"
	"go_refine_dashboard_be/internal/usecases/post/dto"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound             = errors.New("không tìm thấy bài viết")
	ErrPostInvalidInput         = errors.New("dữ liệu bài viết không hợp lệ")
	ErrPostSlugConflict         = errors.New("slug đã tồn tại, vui lòng chọn slug khác")
	ErrPostCategoryNotFound     = errors.New("danh mục bài viết không tồn tại")
	ErrPostCategoryTypeMismatch = errors.New("danh mục không thuộc cùng loại bài viết")
)

var postSlugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

type PostService interface {
	CreatePost(ctx context.Context, authorID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(ctx context.Context, id uint) (*dto.PostResponse, error)
	GetPosts(ctx context.Context, query dto.GetPostsQuery) (*dto.PaginatedPostResponse, error)
	UpdatePost(ctx context.Context, id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(ctx context.Context, id uint) error
}

type postService struct {
	postRepo        repository.PostRepository
	categoryRepo    repository.PostCategoryRepository
	redisClient     *redis.Client
	publicPageCache cache.PublicPageCacheInvalidator
}

func NewPostService(
	postRepo repository.PostRepository,
	categoryRepo repository.PostCategoryRepository,
	redisClient *redis.Client,
	publicPageCache cache.PublicPageCacheInvalidator,
) PostService {
	return &postService{
		postRepo:        postRepo,
		categoryRepo:    categoryRepo,
		redisClient:     redisClient,
		publicPageCache: publicPageCache,
	}
}

func mapPostToResponse(post *entity.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:         post.ID,
		TypeCode:   post.TypeCode,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		AuthorID:   post.AuthorID,
		CategoryID: post.CategoryID,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

func (s *postService) invalidateCache(ctx context.Context, postID uint, typeCode string) {
	if s.redisClient == nil {
		return
	}

	patterns := []string{
		fmt.Sprintf("admin:posts:detail:%d", postID),
		"admin:posts:list:*",
		"posts:list:*",
		fmt.Sprintf("admin:posts:list:type:%s:*", typeCode),
		fmt.Sprintf("posts:type:%s:*", typeCode),
	}

	for _, pattern := range patterns {
		iterator := s.redisClient.Scan(ctx, 0, pattern, 100).Iterator()
		for iterator.Next(ctx) {
			if err := s.redisClient.Del(ctx, iterator.Val()).Err(); err != nil {
				log.Printf("post cache invalidation failed for key %q: %v", iterator.Val(), err)
			}
		}
		if err := iterator.Err(); err != nil {
			log.Printf("post cache scan failed for pattern %q: %v", pattern, err)
		}
	}
}

func (s *postService) CreatePost(ctx context.Context, authorID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	exists, err := s.postRepo.CheckSlugExists(ctx, req.Slug)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra slug: %v", err)
	}
	if exists {
		return nil, ErrPostSlugConflict
	}

	post := &entity.Post{
		TypeCode:   req.TypeCode,
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		AuthorID:   authorID,
		CategoryID: req.CategoryID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("không thể tạo bài viết: %v", err)
	}

	s.invalidateCache(ctx, post.ID, post.TypeCode)

	return mapPostToResponse(post), nil
}

func (s *postService) GetPostByID(ctx context.Context, id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("lỗi lấy chi tiết bài viết: %w", err)
	}

	return mapPostToResponse(post), nil
}

func (s *postService) UpdatePost(ctx context.Context, id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("lỗi lấy bài viết cần cập nhật: %w", err)
	}

	title := strings.TrimSpace(req.Title)
	slug := strings.TrimSpace(req.Slug)
	if title == "" || strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("%w: tiêu đề và nội dung không được để trống", ErrPostInvalidInput)
	}
	if !postSlugPattern.MatchString(slug) {
		return nil, fmt.Errorf("%w: slug chỉ chứa chữ thường, số và dấu gạch ngang", ErrPostInvalidInput)
	}

	exists, err := s.postRepo.CheckSlugExistsExceptID(ctx, slug, id)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra slug: %w", err)
	}
	if exists {
		return nil, ErrPostSlugConflict
	}

	if req.CategoryID != nil {
		category, categoryErr := s.categoryRepo.FindByID(ctx, *req.CategoryID)
		if categoryErr != nil {
			if errors.Is(categoryErr, gorm.ErrRecordNotFound) {
				return nil, ErrPostCategoryNotFound
			}
			return nil, fmt.Errorf("lỗi kiểm tra danh mục bài viết: %w", categoryErr)
		}
		if category.TypeCode != post.TypeCode {
			return nil, ErrPostCategoryTypeMismatch
		}
	}

	post.Title = title
	post.Slug = slug
	post.Content = req.Content
	post.CategoryID = req.CategoryID

	if err := s.postRepo.Update(ctx, post); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") && strings.Contains(strings.ToLower(err.Error()), "slug") {
			return nil, ErrPostSlugConflict
		}
		return nil, fmt.Errorf("không thể cập nhật bài viết: %w", err)
	}

	s.invalidateCache(ctx, post.ID, post.TypeCode)
	s.publicPageCache.InvalidateSource(ctx, pageEntity.PageSectionItemTypePost, post.ID)
	return mapPostToResponse(post), nil
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
		response := mapPostToResponse(&p)
		response.Content = "" // Content is omitted in list.
		data = append(data, *response)
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
	s.publicPageCache.InvalidateSource(ctx, pageEntity.PageSectionItemTypePost, id)

	s.invalidateCache(ctx, id, "")
	return nil
}
