package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/postcategory/entity"

	"gorm.io/gorm"
)

type PostCategoryRepository interface {
	FindTreeByTypeCode(ctx context.Context, typeCode string) ([]entity.PostCategory, error)
	FindBySlug(ctx context.Context, slug string) (*entity.PostCategory, error)
	FindByID(ctx context.Context, id uint) (*entity.PostCategory, error)
	Create(ctx context.Context, category *entity.PostCategory) error
}

type postCategoryRepository struct {
	db *gorm.DB
}

func NewPostCategoryRepository(db *gorm.DB) PostCategoryRepository {
	return &postCategoryRepository{db: db}
}

func (r *postCategoryRepository) FindTreeByTypeCode(ctx context.Context, typeCode string) ([]entity.PostCategory, error) {
	var items []entity.PostCategory
	err := r.db.WithContext(ctx).Where("type_code = ?", typeCode).Order("sort_order ASC, created_at DESC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postCategoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.PostCategory, error) {
	var item entity.PostCategory
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *postCategoryRepository) FindByID(ctx context.Context, id uint) (*entity.PostCategory, error) {
	var item entity.PostCategory
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *postCategoryRepository) Create(ctx context.Context, category *entity.PostCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}
