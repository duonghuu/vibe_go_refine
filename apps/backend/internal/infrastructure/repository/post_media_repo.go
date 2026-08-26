package repository

import (
	"context"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"go_refine_dashboard_be/internal/domain/post/entity"

	"gorm.io/gorm"
)

type PostMediaRepository interface {
	FindByPostID(ctx context.Context, postID uint, collection string) ([]entity.PostMedia, int64, error)
	Create(ctx context.Context, postMedia *entity.PostMedia) error
	CreateAndAttach(ctx context.Context, postMedia *entity.PostMedia) error
	Delete(ctx context.Context, postID uint, mediaID uint, collection string) error
	SyncByCollection(ctx context.Context, postID uint, collection string, items []entity.PostMedia) error
	SyncByCollectionAndAttach(ctx context.Context, postID uint, collection string, items []entity.PostMedia, mediaIDs []uint) error
	GetMediaByID(ctx context.Context, mediaID uint) (*mediaEntity.Media, error)
	GetMediaByIDs(ctx context.Context, mediaIDs []uint) ([]mediaEntity.Media, error)
	UpdateMediaStatus(ctx context.Context, mediaIDs []uint, status mediaEntity.MediaStatus) error
	PostExists(ctx context.Context, postID uint) (bool, error)
}

type postMediaRepository struct {
	db *gorm.DB
}

func NewPostMediaRepository(db *gorm.DB) PostMediaRepository {
	return &postMediaRepository{db: db}
}

func (r *postMediaRepository) FindByPostID(ctx context.Context, postID uint, collection string) ([]entity.PostMedia, int64, error) {
	var items []entity.PostMedia
	var total int64

	q := r.db.WithContext(ctx).Model(&entity.PostMedia{}).Preload("Media").Where("post_id = ?", postID)
	if collection != "" {
		q = q.Where("collection = ?", collection)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *postMediaRepository) Create(ctx context.Context, postMedia *entity.PostMedia) error {
	return r.db.WithContext(ctx).Create(postMedia).Error
}

func (r *postMediaRepository) CreateAndAttach(ctx context.Context, postMedia *entity.PostMedia) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if postMedia.Collection == "thumbnail" {
			if err := tx.Unscoped().Where("post_id = ? AND collection = ?", postMedia.PostID, postMedia.Collection).Delete(&entity.PostMedia{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(postMedia).Error; err != nil {
			return err
		}
		return tx.Model(&mediaEntity.Media{}).
			Where("id = ?", postMedia.MediaID).
			Update("status", mediaEntity.MediaStatusAttached).Error
	})
}

func (r *postMediaRepository) Delete(ctx context.Context, postID uint, mediaID uint, collection string) error {
	q := r.db.WithContext(ctx).Where("post_id = ? AND media_id = ?", postID, mediaID)
	if collection != "" {
		q = q.Where("collection = ?", collection)
	}
	return q.Delete(&entity.PostMedia{}).Error
}

func (r *postMediaRepository) SyncByCollection(ctx context.Context, postID uint, collection string, items []entity.PostMedia) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing associations for this post and collection
		if err := tx.Unscoped().Where("post_id = ? AND collection = ?", postID, collection).Delete(&entity.PostMedia{}).Error; err != nil {
			return err
		}

		// Insert new items
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *postMediaRepository) SyncByCollectionAndAttach(ctx context.Context, postID uint, collection string, items []entity.PostMedia, mediaIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("post_id = ? AND collection = ?", postID, collection).Delete(&entity.PostMedia{}).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		if len(mediaIDs) == 0 {
			return nil
		}
		return tx.Model(&mediaEntity.Media{}).
			Where("id IN ?", mediaIDs).
			Update("status", mediaEntity.MediaStatusAttached).Error
	})
}

func (r *postMediaRepository) GetMediaByID(ctx context.Context, mediaID uint) (*mediaEntity.Media, error) {
	var m mediaEntity.Media
	if err := r.db.WithContext(ctx).First(&m, mediaID).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *postMediaRepository) GetMediaByIDs(ctx context.Context, mediaIDs []uint) ([]mediaEntity.Media, error) {
	if len(mediaIDs) == 0 {
		return []mediaEntity.Media{}, nil
	}
	var media []mediaEntity.Media
	if err := r.db.WithContext(ctx).Where("id IN ?", mediaIDs).Find(&media).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (r *postMediaRepository) UpdateMediaStatus(ctx context.Context, mediaIDs []uint, status mediaEntity.MediaStatus) error {
	if len(mediaIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&mediaEntity.Media{}).
		Where("id IN ?", mediaIDs).
		Update("status", status).Error
}

func (r *postMediaRepository) PostExists(ctx context.Context, postID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("posts").Where("id = ? AND deleted_at IS NULL", postID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
