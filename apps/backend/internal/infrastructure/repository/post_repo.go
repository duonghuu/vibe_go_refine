package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/post/entity"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	FindAndCount(ctx context.Context, title string, typeCode string, offset int, limit int, sort string) ([]entity.Post, int64, error)
	Delete(ctx context.Context, id uint) error
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

func (r *postRepository) FindAndCount(ctx context.Context, title string, typeCode string, offset int, limit int, sort string) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Post{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if typeCode != "" {
		query = query.Where("type_code = ?", typeCode)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("id desc")
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Post{}, id).Error
}
