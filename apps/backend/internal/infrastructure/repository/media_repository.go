package repository

import (
	"context"
	"go_refine_dashboard_be/internal/domain/media/entity"
	"strings"

	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	GetByID(ctx context.Context, id uint) (*entity.Media, error)
	Delete(ctx context.Context, id uint) error
	FindAndCount(ctx context.Context, ownerID uint, isAdmin bool, search, mimeType, status string, offset, limit int) ([]entity.Media, int64, error)
}

func (r *mediaRepository) FindAndCount(ctx context.Context, ownerID uint, isAdmin bool, search, mimeType, status string, offset, limit int) ([]entity.Media, int64, error) {
	var media []entity.Media
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Media{})
	if !isAdmin {
		q = q.Where("owner_id = ?", ownerID)
	}
	if search != "" {
		escaped := strings.NewReplacer(`\\`, `\\\\`, "%", `\\%`, "_", `\\_`).Replace(search)
		pattern := "%" + escaped + "%"
		q = q.Where("(file_name LIKE ? OR original_name LIKE ?)", pattern, pattern)
	}
	if mimeType == "image" {
		q = q.Where("mime_type LIKE ?", "image/%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC, id DESC").Offset(offset).Limit(limit).Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
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
