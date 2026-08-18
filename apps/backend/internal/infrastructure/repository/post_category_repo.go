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
	FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, status, typeCode string) ([]entity.PostCategory, int64, error)
	Delete(ctx context.Context, id uint) error
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

func (r *postCategoryRepository) FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, status, typeCode string) ([]entity.PostCategory, int64, error) {
	var items []entity.PostCategory
	var total int64
	dbQuery := r.db.WithContext(ctx).Model(&entity.PostCategory{})

	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ? OR slug LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}
	if typeCode != "" {
		dbQuery = dbQuery.Where("type_code = ?", typeCode)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortField != "" {
		dbQuery = dbQuery.Order(sortField + " " + sortOrder)
	} else {
		dbQuery = dbQuery.Order("sort_order ASC, created_at DESC")
	}

	if limit > 0 {
		dbQuery = dbQuery.Offset(skip).Limit(limit)
	}

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *postCategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.PostCategory{}, id).Error
}
