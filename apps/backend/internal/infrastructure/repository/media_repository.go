package repository

import (
	"context"
	"go_refine_dashboard_be/internal/domain/media/entity"

	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	GetByID(ctx context.Context, id uint) (*entity.Media, error)
	Delete(ctx context.Context, id uint) error
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(ctx context.Context, media *entity.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *mediaRepository) GetByID(ctx context.Context, id uint) (*entity.Media, error) {
	var media entity.Media
	err := r.db.WithContext(ctx).First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) Delete(ctx context.Context, id uint) error {
	// Logic: Hard Delete for physical cleanup, but model has DeletedAt.
	// Since the plan says "xóa file vật lý tương ứng, tiến hành Hard Delete bản ghi khỏi DB", we use Unscoped.
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.Media{}, id).Error
}
