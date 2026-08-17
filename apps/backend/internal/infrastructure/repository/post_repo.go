package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/post/entity"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Post{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
