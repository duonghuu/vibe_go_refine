package repository

import (
	"context"

	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"

	"gorm.io/gorm"
)

// PublicPageCacheRepository resolves the public page slugs affected by a
// source mutation. It intentionally returns slugs instead of Redis keys so
// cache-key construction stays in one place.
type PublicPageCacheRepository interface {
	FindPageSlug(ctx context.Context, pageID uint) (string, error)
	FindPageSlugsBySource(ctx context.Context, itemType pageEntity.PageSectionItemType, itemIDs []uint) ([]string, error)
	FindPageSlugsByMedia(ctx context.Context, mediaIDs []uint) ([]string, error)
}

type publicPageCacheRepository struct {
	db *gorm.DB
}

func NewPublicPageCacheRepository(db *gorm.DB) PublicPageCacheRepository {
	return &publicPageCacheRepository{db: db}
}

func (r *publicPageCacheRepository) FindPageSlug(ctx context.Context, pageID uint) (string, error) {
	var row struct {
		Slug string
	}
	if err := r.db.WithContext(ctx).
		Table("pages").
		Select("slug").
		Where("id = ?", pageID).
		Limit(1).
		Scan(&row).Error; err != nil {
		return "", err
	}
	return row.Slug, nil
}

func (r *publicPageCacheRepository) FindPageSlugsBySource(ctx context.Context, itemType pageEntity.PageSectionItemType, itemIDs []uint) ([]string, error) {
	if len(itemIDs) == 0 {
		return []string{}, nil
	}

	var slugs []string
	err := r.db.WithContext(ctx).
		Table("pages AS p").
		Select("DISTINCT p.slug").
		Joins("JOIN page_sections AS ps ON ps.page_id = p.id").
		Joins("JOIN page_section_items AS psi ON psi.section_id = ps.id").
		Where("p.status = ? AND p.deleted_at IS NULL", pageEntity.PageStatusPublished).
		Where("ps.status = ? AND ps.deleted_at IS NULL", pageEntity.PageSectionStatusActive).
		Where("psi.item_type = ? AND psi.item_id IN ? AND psi.deleted_at IS NULL", itemType, itemIDs).
		Where("p.slug <> ''").
		Order("p.slug ASC").
		Pluck("p.slug", &slugs).Error
	if err != nil {
		return nil, err
	}
	return slugs, nil
}

func (r *publicPageCacheRepository) FindPageSlugsByMedia(ctx context.Context, mediaIDs []uint) ([]string, error) {
	if len(mediaIDs) == 0 {
		return []string{}, nil
	}

	var slugs []string
	err := r.db.WithContext(ctx).
		Table("pages AS p").
		Select("DISTINCT p.slug").
		Joins("LEFT JOIN page_media AS pm ON pm.page_id = p.id AND pm.deleted_at IS NULL").
		Joins("LEFT JOIN page_sections AS ps ON ps.page_id = p.id AND ps.deleted_at IS NULL AND ps.status = ?", pageEntity.PageSectionStatusActive).
		Joins("LEFT JOIN page_section_items AS psi ON psi.section_id = ps.id AND psi.deleted_at IS NULL AND psi.item_type = ?", pageEntity.PageSectionItemTypeMedia).
		Where("p.status = ? AND p.deleted_at IS NULL", pageEntity.PageStatusPublished).
		Where("(pm.media_id IN ? OR ps.background_media_id IN ? OR ps.feature_media_id IN ? OR psi.item_id IN ?)", mediaIDs, mediaIDs, mediaIDs, mediaIDs).
		Where("p.slug <> ''").
		Order("p.slug ASC").
		Pluck("p.slug", &slugs).Error
	if err != nil {
		return nil, err
	}
	return slugs, nil
}
