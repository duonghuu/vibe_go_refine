package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/category/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string, parentID *uint) ([]entity.Category, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint) error
	BulkUpdateStatus(ctx context.Context, ids []uint, status string) error
	Reorder(ctx context.Context, items []struct {
		ID        uint
		SortOrder int
	}) error
	CountProductsByCategoryID(ctx context.Context, categoryID uint) (int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string, parentID *uint) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	q := r.db.WithContext(ctx).Model(&entity.Category{})

	if query != "" {
		q = q.Where("name LIKE ? OR slug LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if parentID != nil {
		q = q.Where("parent_id = ?", *parentID)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortField != "" {
		q = q.Order(sortField + " " + sortOrder)
	} else {
		q = q.Order("sort_order ASC, created_at DESC")
	}

	if limit > 0 {
		q = q.Offset(skip).Limit(limit)
	}

	if err := q.Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uint) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

func (r *categoryRepository) BulkUpdateStatus(ctx context.Context, ids []uint, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Category{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *categoryRepository) Reorder(ctx context.Context, items []struct {
	ID        uint
	SortOrder int
}) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&entity.Category{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *categoryRepository) CountProductsByCategoryID(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	// Giả định bảng products có trường category_id. Vì model Product chưa chắc có,
	// chúng ta dùng Table("products") để query raw count nhằm tránh lỗi dependency vòng nếu entity chưa chuẩn.
	err := r.db.WithContext(ctx).Table("products").Where("category_id = ? AND deleted_at IS NULL", categoryID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
