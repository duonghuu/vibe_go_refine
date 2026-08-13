package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/posttype/entity"

	"gorm.io/gorm"
)

type PostTypeRepository interface {
	FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string) ([]entity.PostType, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.PostType, error)
	FindByCode(ctx context.Context, code string) (*entity.PostType, error)
	Create(ctx context.Context, postType *entity.PostType) error
	Update(ctx context.Context, postType *entity.PostType) error
	Delete(ctx context.Context, id uint) error
	CountPostsByPostTypeID(ctx context.Context, id uint) (int64, error)
}

type postTypeRepository struct {
	db *gorm.DB
}

func NewPostTypeRepository(db *gorm.DB) PostTypeRepository {
	return &postTypeRepository{db: db}
}

func (r *postTypeRepository) FindAll(ctx context.Context, skip, limit int, sortField, sortOrder, query, status string) ([]entity.PostType, int64, error) {
	var items []entity.PostType
	var total int64

	q := r.db.WithContext(ctx).Model(&entity.PostType{})

	if query != "" {
		q = q.Where("name LIKE ? OR code LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
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

	if err := q.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *postTypeRepository) FindByID(ctx context.Context, id uint) (*entity.PostType, error) {
	var item entity.PostType
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *postTypeRepository) FindByCode(ctx context.Context, code string) (*entity.PostType, error) {
	var item entity.PostType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *postTypeRepository) Create(ctx context.Context, postType *entity.PostType) error {
	return r.db.WithContext(ctx).Create(postType).Error
}

func (r *postTypeRepository) Update(ctx context.Context, postType *entity.PostType) error {
	return r.db.WithContext(ctx).Save(postType).Error
}

func (r *postTypeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.PostType{}, id).Error
}

func (r *postTypeRepository) CountPostsByPostTypeID(ctx context.Context, id uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("posts").Where("post_type_id = ? AND deleted_at IS NULL", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
