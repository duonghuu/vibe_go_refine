package repository

import (
	"context"
	"errors"
	"time"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	seoEntity "go_refine_dashboard_be/internal/domain/seometa/entity"

	"gorm.io/gorm"
)

type PublicPageProjection struct {
	ID        uint
	Title     string
	Slug      string
	Content   string
	UpdatedAt time.Time
	SEO       *seoEntity.SEOMeta
	Thumbnail *PublicPageMediaProjection
	Gallery   []PublicPageMediaProjection
	Sections  []PublicPageSectionProjection
}

type PublicPageMediaProjection struct {
	ID           uint
	OriginalURL  string
	ThumbnailURL string
	MediumURL    string
	SortOrder    int
}

type PublicPageSectionProjection struct {
	ID              uint
	Key             string
	Title           string
	Description     string
	BackgroundColor *string
	BackgroundMedia *PublicSectionMediaProjection
	FeatureMedia    *PublicSectionMediaProjection
	SortOrder       int
	Collections     map[string][]PublicPageSectionItemProjection
}

type PublicSectionMediaProjection struct {
	ID           uint
	OriginalURL  string
	ThumbnailURL string
	MediumURL    string
}

type PublicPageSectionItemProjection struct {
	ItemType  string
	ItemID    uint
	SortOrder int
	Data      PublicSectionItemDataProjection
}

type PublicSectionItemDataProjection struct {
	ID           uint
	Name         string
	Description  string
	Title        string
	Slug         string
	TypeCode     string
	ImageURL     string
	OriginalURL  string
	ThumbnailURL string
	MediumURL    string
}

type PublicPageRepository interface {
	FindPublishedBySlug(ctx context.Context, slug string) (*PublicPageProjection, error)
}

type publicPageRepository struct {
	db *gorm.DB
}

func NewPublicPageRepository(db *gorm.DB) PublicPageRepository {
	return &publicPageRepository{db: db}
}

func (r *publicPageRepository) FindPublishedBySlug(ctx context.Context, slug string) (*PublicPageProjection, error) {
	var page struct {
		ID        uint
		Title     string
		Slug      string
		Content   string
		UpdatedAt time.Time
	}
	if err := r.db.WithContext(ctx).
		Table("pages").
		Select("id, title, slug, content, updated_at").
		Where("slug = ? AND status = ? AND deleted_at IS NULL", slug, pageEntity.PageStatusPublished).
		First(&page).Error; err != nil {
		return nil, err
	}

	projection := &PublicPageProjection{
		ID:        page.ID,
		Title:     page.Title,
		Slug:      page.Slug,
		Content:   page.Content,
		UpdatedAt: page.UpdatedAt,
		Gallery:   make([]PublicPageMediaProjection, 0),
		Sections:  make([]PublicPageSectionProjection, 0),
	}

	var seo seoEntity.SEOMeta
	seoErr := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", "page", page.ID).
		First(&seo).Error
	if seoErr == nil {
		projection.SEO = &seo
	} else if !errors.Is(seoErr, gorm.ErrRecordNotFound) {
		return nil, seoErr
	}

	var mediaRows []struct {
		LinkID       uint
		MediaID      uint
		Collection   string
		SortOrder    int
		OriginalURL  string
		ThumbnailURL string
		MediumURL    string
	}
	if err := r.db.WithContext(ctx).
		Table("page_media AS pm").
		Select("pm.id AS link_id, m.id AS media_id, pm.collection, pm.sort_order, m.original_url, m.thumbnail_url, m.medium_url").
		Joins("JOIN media AS m ON m.id = pm.media_id").
		Where("pm.page_id = ? AND pm.collection IN ? AND pm.deleted_at IS NULL AND m.deleted_at IS NULL AND m.status = ?", page.ID, []string{string(pageEntity.PageMediaCollectionThumbnail), string(pageEntity.PageMediaCollectionGallery)}, mediaEntity.MediaStatusAttached).
		Order("pm.collection ASC, pm.sort_order ASC, pm.id ASC").
		Find(&mediaRows).Error; err != nil {
		return nil, err
	}
	for _, row := range mediaRows {
		media := PublicPageMediaProjection{ID: row.MediaID, OriginalURL: row.OriginalURL, ThumbnailURL: row.ThumbnailURL, MediumURL: row.MediumURL, SortOrder: row.SortOrder}
		if row.Collection == string(pageEntity.PageMediaCollectionThumbnail) && projection.Thumbnail == nil {
			projection.Thumbnail = &media
			continue
		}
		if row.Collection == string(pageEntity.PageMediaCollectionGallery) {
			projection.Gallery = append(projection.Gallery, media)
		}
	}

	var sectionRows []struct {
		ID                uint
		Key               string
		Title             string
		Description       string
		BackgroundColor   *string
		BackgroundMediaID *uint
		FeatureMediaID    *uint
		SortOrder         int
	}
	if err := r.db.WithContext(ctx).
		Table("page_sections").
		Select("id, `key`, title, description, background_color, background_media_id, feature_media_id, sort_order").
		Where("page_id = ? AND status = ? AND deleted_at IS NULL", page.ID, pageEntity.PageSectionStatusActive).
		Order("sort_order ASC, id ASC").
		Find(&sectionRows).Error; err != nil {
		return nil, err
	}
	if len(sectionRows) == 0 {
		return projection, nil
	}

	sectionIDs := make([]uint, 0, len(sectionRows))
	mediaIDs := make([]uint, 0, len(sectionRows)*2)
	for _, row := range sectionRows {
		sectionIDs = append(sectionIDs, row.ID)
		if row.BackgroundMediaID != nil {
			mediaIDs = append(mediaIDs, *row.BackgroundMediaID)
		}
		if row.FeatureMediaID != nil {
			mediaIDs = append(mediaIDs, *row.FeatureMediaID)
		}
	}

	sectionMedia, err := r.findSectionMedia(ctx, mediaIDs)
	if err != nil {
		return nil, err
	}

	var itemRows []struct {
		SectionID  uint
		ItemType   string
		ItemID     uint
		Collection string
		SortOrder  int
	}
	if err := r.db.WithContext(ctx).
		Table("page_section_items").
		Select("section_id, item_type, item_id, collection, sort_order").
		Where("section_id IN ? AND deleted_at IS NULL", sectionIDs).
		Order("section_id ASC, collection ASC, sort_order ASC, id ASC").
		Find(&itemRows).Error; err != nil {
		return nil, err
	}

	sourceData, err := r.findPublicSectionSources(ctx, itemRows)
	if err != nil {
		return nil, err
	}
	sectionMap := make(map[uint]*PublicPageSectionProjection, len(sectionRows))
	for _, row := range sectionRows {
		section := PublicPageSectionProjection{
			ID: row.ID, Key: row.Key, Title: row.Title, Description: row.Description,
			BackgroundColor: row.BackgroundColor, SortOrder: row.SortOrder,
			Collections: make(map[string][]PublicPageSectionItemProjection),
		}
		if row.BackgroundMediaID != nil {
			section.BackgroundMedia = sectionMedia[*row.BackgroundMediaID]
		}
		if row.FeatureMediaID != nil {
			section.FeatureMedia = sectionMedia[*row.FeatureMediaID]
		}
		projection.Sections = append(projection.Sections, section)
		sectionMap[row.ID] = &projection.Sections[len(projection.Sections)-1]
	}
	for _, row := range itemRows {
		data, ok := sourceData[row.ItemType][row.ItemID]
		if !ok {
			continue
		}
		section := sectionMap[row.SectionID]
		section.Collections[row.Collection] = append(section.Collections[row.Collection], PublicPageSectionItemProjection{
			ItemType: row.ItemType, ItemID: row.ItemID, SortOrder: row.SortOrder, Data: data,
		})
	}

	return projection, nil
}

func (r *publicPageRepository) findSectionMedia(ctx context.Context, ids []uint) (map[uint]*PublicSectionMediaProjection, error) {
	result := make(map[uint]*PublicSectionMediaProjection)
	if len(ids) == 0 {
		return result, nil
	}
	var rows []PublicSectionMediaProjection
	if err := r.db.WithContext(ctx).
		Table("media").
		Select("id, original_url, thumbnail_url, medium_url").
		Where("id IN ? AND deleted_at IS NULL AND status = ? AND LOWER(mime_type) LIKE ?", ids, mediaEntity.MediaStatusAttached, "image/%").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		result[rows[i].ID] = &rows[i]
	}
	return result, nil
}

func (r *publicPageRepository) findPublicSectionSources(ctx context.Context, rows []struct {
	SectionID  uint
	ItemType   string
	ItemID     uint
	Collection string
	SortOrder  int
}) (map[string]map[uint]PublicSectionItemDataProjection, error) {
	idsByType := make(map[string][]uint)
	for _, row := range rows {
		idsByType[row.ItemType] = append(idsByType[row.ItemType], row.ItemID)
	}
	result := map[string]map[uint]PublicSectionItemDataProjection{
		string(pageEntity.PageSectionItemTypeCategory): {},
		string(pageEntity.PageSectionItemTypePost):     {},
		string(pageEntity.PageSectionItemTypeMedia):    {},
	}

	if ids := uniqueUintIDs(idsByType[string(pageEntity.PageSectionItemTypeCategory)]); len(ids) > 0 {
		var source []PublicSectionItemDataProjection
		if err := r.db.WithContext(ctx).Table("categories").Select("id, name, description, slug, image_url").Where("id IN ? AND status = ? AND deleted_at IS NULL", ids, "ACTIVE").Find(&source).Error; err != nil {
			return nil, err
		}
		for _, item := range source {
			result[string(pageEntity.PageSectionItemTypeCategory)][item.ID] = item
		}
	}
	if ids := uniqueUintIDs(idsByType[string(pageEntity.PageSectionItemTypePost)]); len(ids) > 0 {
		var source []PublicSectionItemDataProjection
		if err := r.db.WithContext(ctx).
			Table("posts AS p").
			Select("p.id, p.title, p.content AS description, p.slug, p.type_code, m.original_url, m.thumbnail_url, m.medium_url").
			Joins("LEFT JOIN post_media AS pm ON pm.post_id = p.id AND pm.collection = ? AND pm.deleted_at IS NULL", "thumbnail").
			Joins("LEFT JOIN media AS m ON m.id = pm.media_id AND m.deleted_at IS NULL AND m.status = ? AND LOWER(m.mime_type) LIKE ?", mediaEntity.MediaStatusAttached, "image/%").
			Where("p.id IN ? AND p.deleted_at IS NULL", ids).
			Order("p.id ASC, pm.sort_order ASC, pm.id ASC").
			Find(&source).Error; err != nil {
			return nil, err
		}
		for _, item := range source {
			if _, exists := result[string(pageEntity.PageSectionItemTypePost)][item.ID]; !exists {
				result[string(pageEntity.PageSectionItemTypePost)][item.ID] = item
			}
		}
	}
	if ids := uniqueUintIDs(idsByType[string(pageEntity.PageSectionItemTypeMedia)]); len(ids) > 0 {
		var source []PublicSectionItemDataProjection
		if err := r.db.WithContext(ctx).Table("media").Select("id, original_url, thumbnail_url, medium_url").Where("id IN ? AND deleted_at IS NULL AND status = ? AND LOWER(mime_type) LIKE ?", ids, mediaEntity.MediaStatusAttached, "image/%").Find(&source).Error; err != nil {
			return nil, err
		}
		for _, item := range source {
			result[string(pageEntity.PageSectionItemTypeMedia)][item.ID] = item
		}
	}
	return result, nil
}

func uniqueUintIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
