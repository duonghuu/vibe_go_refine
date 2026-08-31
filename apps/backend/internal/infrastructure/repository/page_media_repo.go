package repository

import (
	"context"
	"database/sql"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageMediaRepository interface {
	FindPageForUpdate(ctx context.Context, tx *gorm.DB, pageID uint) (*pageEntity.Page, error)
	FindByPageID(ctx context.Context, pageID uint, collection string) ([]pageEntity.PageMedia, int64, error)
	FindMediaByIDsForUpdate(ctx context.Context, tx *gorm.DB, ids []uint) ([]mediaEntity.Media, error)
	FindActiveLinks(ctx context.Context, tx *gorm.DB, pageID, mediaID uint, collection string) ([]pageEntity.PageMedia, error)
	CreateOrRestore(ctx context.Context, tx *gorm.DB, item *pageEntity.PageMedia) error
	SoftDeleteNotIn(ctx context.Context, tx *gorm.DB, pageID uint, collection string, mediaIDs []uint) error
	SoftDeleteLink(ctx context.Context, tx *gorm.DB, pageID, mediaID uint, collection string) (bool, error)
	MarkAttached(ctx context.Context, tx *gorm.DB, mediaIDs []uint) error
	NextSortOrder(ctx context.Context, tx *gorm.DB, pageID uint, collection string) (int, error)
}

type pageMediaRepository struct{ db *gorm.DB }

func NewPageMediaRepository(db *gorm.DB) PageMediaRepository { return &pageMediaRepository{db: db} }

func (r *pageMediaRepository) FindPageForUpdate(ctx context.Context, tx *gorm.DB, pageID uint) (*pageEntity.Page, error) {
	var page pageEntity.Page
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&page, pageID).Error
	return &page, err
}

func (r *pageMediaRepository) FindByPageID(ctx context.Context, pageID uint, collection string) ([]pageEntity.PageMedia, int64, error) {
	var items []pageEntity.PageMedia
	q := r.db.WithContext(ctx).Model(&pageEntity.PageMedia{}).Preload("Media").Where("page_id = ?", pageID)
	if collection != "" {
		q = q.Where("collection = ?", collection)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *pageMediaRepository) FindMediaByIDsForUpdate(ctx context.Context, tx *gorm.DB, ids []uint) ([]mediaEntity.Media, error) {
	if len(ids) == 0 {
		return []mediaEntity.Media{}, nil
	}
	var media []mediaEntity.Media
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Find(&media).Error
	return media, err
}

func (r *pageMediaRepository) FindActiveLinks(ctx context.Context, tx *gorm.DB, pageID, mediaID uint, collection string) ([]pageEntity.PageMedia, error) {
	var items []pageEntity.PageMedia
	q := tx.WithContext(ctx).Where("page_id = ? AND media_id = ?", pageID, mediaID)
	if collection != "" {
		q = q.Where("collection = ?", collection)
	}
	err := q.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *pageMediaRepository) CreateOrRestore(ctx context.Context, tx *gorm.DB, item *pageEntity.PageMedia) error {
	var existing pageEntity.PageMedia
	err := tx.WithContext(ctx).Unscoped().Where("page_id = ? AND media_id = ? AND collection = ?", item.PageID, item.MediaID, item.Collection).First(&existing).Error
	if err == nil {
		return tx.WithContext(ctx).Unscoped().Model(&existing).Updates(map[string]interface{}{"sort_order": item.SortOrder, "deleted_at": nil}).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return tx.WithContext(ctx).Create(item).Error
}

func (r *pageMediaRepository) SoftDeleteNotIn(ctx context.Context, tx *gorm.DB, pageID uint, collection string, mediaIDs []uint) error {
	q := tx.WithContext(ctx).Where("page_id = ? AND collection = ?", pageID, collection)
	if len(mediaIDs) > 0 {
		q = q.Where("media_id NOT IN ?", mediaIDs)
	}
	return q.Delete(&pageEntity.PageMedia{}).Error
}

func (r *pageMediaRepository) SoftDeleteLink(ctx context.Context, tx *gorm.DB, pageID, mediaID uint, collection string) (bool, error) {
	result := tx.WithContext(ctx).Where("page_id = ? AND media_id = ? AND collection = ?", pageID, mediaID, collection).Delete(&pageEntity.PageMedia{})
	return result.RowsAffected > 0, result.Error
}

func (r *pageMediaRepository) MarkAttached(ctx context.Context, tx *gorm.DB, mediaIDs []uint) error {
	if len(mediaIDs) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Model(&mediaEntity.Media{}).Where("id IN ?", mediaIDs).Update("status", mediaEntity.MediaStatusAttached).Error
}

func (r *pageMediaRepository) NextSortOrder(ctx context.Context, tx *gorm.DB, pageID uint, collection string) (int, error) {
	var max sql.NullInt64
	err := tx.WithContext(ctx).Model(&pageEntity.PageMedia{}).Where("page_id = ? AND collection = ?", pageID, collection).Select("MAX(sort_order)").Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if !max.Valid {
		return 0, nil
	}
	return int(max.Int64) + 1, nil
}
