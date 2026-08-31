package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/page/entity"

	"gorm.io/gorm"
)

type PageRepository interface {
	Create(ctx context.Context, page *entity.Page) error
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
}

type pageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &pageRepository{db: db}
}

func (r *pageRepository) Create(ctx context.Context, page *entity.Page) error {
	return r.db.WithContext(ctx).Create(page).Error
}

func (r *pageRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Page{}).
		Where("slug = ?", slug).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
