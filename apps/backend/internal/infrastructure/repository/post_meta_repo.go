package repository

import (
	"context"

	"go_refine_dashboard_be/internal/domain/post/entity"

	"gorm.io/gorm"
)

type PostMetaRepository interface {
	FindByPostID(ctx context.Context, postID uint) ([]entity.PostMeta, error)
	Sync(ctx context.Context, postID uint, items []entity.PostMeta) error
	Delete(ctx context.Context, postID uint, key string) error
}

type postMetaRepository struct {
	db *gorm.DB
}

func NewPostMetaRepository(db *gorm.DB) PostMetaRepository {
	return &postMetaRepository{db: db}
}

func (r *postMetaRepository) FindByPostID(ctx context.Context, postID uint) ([]entity.PostMeta, error) {
	var items []entity.PostMeta
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *postMetaRepository) Sync(ctx context.Context, postID uint, items []entity.PostMeta) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", postID).Delete(&entity.PostMeta{}).Error; err != nil {
			return err
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *postMetaRepository) Delete(ctx context.Context, postID uint, key string) error {
	return r.db.WithContext(ctx).Where("post_id = ? AND key = ?", postID, key).Delete(&entity.PostMeta{}).Error
}
