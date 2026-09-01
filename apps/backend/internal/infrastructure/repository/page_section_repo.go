package repository

import (
	"context"
	"fmt"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	pageEntity "go_refine_dashboard_be/internal/domain/page/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SectionOrderUpdate struct {
	ID        uint
	SortOrder int
}
type PageSectionRepository interface {
	FindPageForUpdate(context.Context, *gorm.DB, uint) (*pageEntity.Page, error)
	FindSectionForUpdate(context.Context, *gorm.DB, uint, uint) (*pageEntity.PageSection, error)
	FindActiveSectionsForUpdate(context.Context, *gorm.DB, uint) ([]pageEntity.PageSection, error)
	FindMediaByIDsForUpdate(context.Context, *gorm.DB, []uint) ([]mediaEntity.Media, error)
	FindSourceIDsForUpdate(context.Context, *gorm.DB, pageEntity.PageSectionItemType, []uint) ([]uint, error)
	CreateOrRestoreSection(context.Context, *gorm.DB, *pageEntity.PageSection) error
	UpdateSection(context.Context, *gorm.DB, *pageEntity.PageSection) error
	SoftDeleteSectionAndItems(context.Context, *gorm.DB, uint) error
	BulkUpdateSectionOrder(context.Context, *gorm.DB, uint, []SectionOrderUpdate) error
	FindCollectionItemsForUpdate(context.Context, *gorm.DB, uint, string) ([]pageEntity.PageSectionItem, error)
	CreateOrRestoreItem(context.Context, *gorm.DB, *pageEntity.PageSectionItem) error
	SoftDeleteItemsNotIn(context.Context, *gorm.DB, uint, string, pageEntity.PageSectionItemType, []uint) error
	SoftDeleteCollection(context.Context, *gorm.DB, uint, string) error
	MarkMediaAttached(context.Context, *gorm.DB, []uint) error
	FindSections(context.Context, uint, string, int, int) ([]pageEntity.PageSection, int64, error)
	FindSectionItems(context.Context, uint, uint, string, int, int) ([]pageEntity.PageSectionItem, int64, error)
	FindCollectionSummaries(context.Context, []uint) (map[uint][]SectionSummary, error)
	FindSourceData(context.Context, pageEntity.PageSectionItemType, []uint) (map[uint]SourceData, error)
}
type SectionSummary struct {
	SectionID  uint
	Collection string
	ItemType   string
	Total      int64
}
type SourceData struct {
	ID                                                                                                      uint
	Name, Title, Slug, TypeCode, ImageURL, Status, FileName, OriginalURL, ThumbnailURL, MediumURL, MimeType string
}
type pageSectionRepository struct{ db *gorm.DB }

func NewPageSectionRepository(db *gorm.DB) PageSectionRepository {
	return &pageSectionRepository{db: db}
}
func (r *pageSectionRepository) FindPageForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*pageEntity.Page, error) {
	var v pageEntity.Page
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&v, id).Error
	return &v, err
}
func (r *pageSectionRepository) FindSectionForUpdate(ctx context.Context, tx *gorm.DB, pageID, id uint) (*pageEntity.PageSection, error) {
	var v pageEntity.PageSection
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("page_id = ? AND id = ?", pageID, id).First(&v).Error
	return &v, err
}
func (r *pageSectionRepository) FindActiveSectionsForUpdate(ctx context.Context, tx *gorm.DB, pageID uint) ([]pageEntity.PageSection, error) {
	var v []pageEntity.PageSection
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("page_id = ?", pageID).Order("sort_order ASC, id ASC").Find(&v).Error
	return v, err
}
func (r *pageSectionRepository) FindMediaByIDsForUpdate(ctx context.Context, tx *gorm.DB, ids []uint) ([]mediaEntity.Media, error) {
	if len(ids) == 0 {
		return []mediaEntity.Media{}, nil
	}
	var v []mediaEntity.Media
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Find(&v).Error
	return v, err
}
func (r *pageSectionRepository) FindSourceIDsForUpdate(ctx context.Context, tx *gorm.DB, typ pageEntity.PageSectionItemType, ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return []uint{}, nil
	}
	var out []uint
	var table string
	switch typ {
	case pageEntity.PageSectionItemTypeCategory:
		table = "categories"
	case pageEntity.PageSectionItemTypePost:
		table = "posts"
	case pageEntity.PageSectionItemTypeMedia:
		table = "media"
	default:
		return nil, fmt.Errorf("invalid source type")
	}
	err := tx.WithContext(ctx).Table(table).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ? AND deleted_at IS NULL", ids).Pluck("id", &out).Error
	return out, err
}
func (r *pageSectionRepository) CreateOrRestoreSection(ctx context.Context, tx *gorm.DB, v *pageEntity.PageSection) error {
	var old pageEntity.PageSection
	err := tx.WithContext(ctx).Unscoped().Where("page_id = ? AND `key` = ?", v.PageID, v.Key).First(&old).Error
	if err == nil {
		if !old.DeletedAt.Valid {
			return gorm.ErrDuplicatedKey
		}
		v.ID = old.ID
		return tx.WithContext(ctx).Unscoped().Model(&old).Updates(map[string]interface{}{"name": v.Name, "title": v.Title, "description": v.Description, "background_color": v.BackgroundColor, "background_media_id": v.BackgroundMediaID, "feature_media_id": v.FeatureMediaID, "sort_order": v.SortOrder, "status": v.Status, "deleted_at": nil}).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return tx.WithContext(ctx).Create(v).Error
}
func (r *pageSectionRepository) UpdateSection(ctx context.Context, tx *gorm.DB, v *pageEntity.PageSection) error {
	return tx.WithContext(ctx).Model(&pageEntity.PageSection{}).Where("id = ? AND page_id = ?", v.ID, v.PageID).Updates(map[string]interface{}{"key": v.Key, "name": v.Name, "title": v.Title, "description": v.Description, "background_color": v.BackgroundColor, "background_media_id": v.BackgroundMediaID, "feature_media_id": v.FeatureMediaID, "status": v.Status}).Error
}
func (r *pageSectionRepository) SoftDeleteSectionAndItems(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := tx.WithContext(ctx).Where("section_id = ?", id).Delete(&pageEntity.PageSectionItem{}).Error; err != nil {
		return err
	}
	return tx.WithContext(ctx).Delete(&pageEntity.PageSection{}, id).Error
}
func (r *pageSectionRepository) BulkUpdateSectionOrder(ctx context.Context, tx *gorm.DB, pageID uint, items []SectionOrderUpdate) error {
	for _, item := range items {
		if err := tx.WithContext(ctx).Model(&pageEntity.PageSection{}).Where("id = ? AND page_id = ?", item.ID, pageID).Update("sort_order", item.SortOrder).Error; err != nil {
			return err
		}
	}
	return nil
}
func (r *pageSectionRepository) FindCollectionItemsForUpdate(ctx context.Context, tx *gorm.DB, id uint, collection string) ([]pageEntity.PageSectionItem, error) {
	var v []pageEntity.PageSectionItem
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("section_id = ? AND collection = ?", id, collection).Order("sort_order ASC, id ASC").Find(&v).Error
	return v, err
}
func (r *pageSectionRepository) CreateOrRestoreItem(ctx context.Context, tx *gorm.DB, v *pageEntity.PageSectionItem) error {
	var old pageEntity.PageSectionItem
	err := tx.WithContext(ctx).Unscoped().Where("section_id = ? AND collection = ? AND item_type = ? AND item_id = ?", v.SectionID, v.Collection, v.ItemType, v.ItemID).First(&old).Error
	if err == nil {
		if !old.DeletedAt.Valid {
			return gorm.ErrDuplicatedKey
		}
		return tx.WithContext(ctx).Unscoped().Model(&old).Updates(map[string]interface{}{"sort_order": v.SortOrder, "deleted_at": nil}).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return tx.WithContext(ctx).Create(v).Error
}
func (r *pageSectionRepository) SoftDeleteItemsNotIn(ctx context.Context, tx *gorm.DB, sectionID uint, collection string, typ pageEntity.PageSectionItemType, ids []uint) error {
	q := tx.WithContext(ctx).Where("section_id = ? AND collection = ? AND item_type = ?", sectionID, collection, typ)
	if len(ids) > 0 {
		q = q.Where("item_id NOT IN ?", ids)
	}
	return q.Delete(&pageEntity.PageSectionItem{}).Error
}
func (r *pageSectionRepository) SoftDeleteCollection(ctx context.Context, tx *gorm.DB, sectionID uint, collection string) error {
	return tx.WithContext(ctx).Where("section_id = ? AND collection = ?", sectionID, collection).Delete(&pageEntity.PageSectionItem{}).Error
}
func (r *pageSectionRepository) MarkMediaAttached(ctx context.Context, tx *gorm.DB, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Model(&mediaEntity.Media{}).Where("id IN ?", ids).Update("status", mediaEntity.MediaStatusAttached).Error
}
func (r *pageSectionRepository) FindSections(ctx context.Context, pageID uint, status string, offset, limit int) ([]pageEntity.PageSection, int64, error) {
	var v []pageEntity.PageSection
	var total int64
	q := r.db.WithContext(ctx).Model(&pageEntity.PageSection{}).Preload("BackgroundMedia").Preload("FeatureMedia").Where("page_id = ?", pageID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("sort_order ASC,id ASC").Offset(offset).Limit(limit).Find(&v).Error
	return v, total, err
}
func (r *pageSectionRepository) FindSectionItems(ctx context.Context, pageID, sectionID uint, collection string, offset, limit int) ([]pageEntity.PageSectionItem, int64, error) {
	var v []pageEntity.PageSectionItem
	var total int64
	q := r.db.WithContext(ctx).Model(&pageEntity.PageSectionItem{}).Joins("JOIN page_sections ON page_sections.id = page_section_items.section_id").Where("page_sections.id = ? AND page_sections.page_id = ? AND page_section_items.collection = ?", sectionID, pageID, collection)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("page_section_items.sort_order ASC,page_section_items.id ASC").Offset(offset).Limit(limit).Find(&v).Error
	return v, total, err
}
func (r *pageSectionRepository) FindCollectionSummaries(ctx context.Context, ids []uint) (map[uint][]SectionSummary, error) {
	out := make(map[uint][]SectionSummary)
	if len(ids) == 0 {
		return out, nil
	}
	var rows []SectionSummary
	err := r.db.WithContext(ctx).Model(&pageEntity.PageSectionItem{}).Select("section_id, collection, item_type, COUNT(*) AS total").Where("section_id IN ?", ids).Group("section_id, collection, item_type").Scan(&rows).Error
	for _, v := range rows {
		out[v.SectionID] = append(out[v.SectionID], v)
	}
	return out, err
}
func (r *pageSectionRepository) FindSourceData(ctx context.Context, typ pageEntity.PageSectionItemType, ids []uint) (map[uint]SourceData, error) {
	out := make(map[uint]SourceData)
	if len(ids) == 0 {
		return out, nil
	}
	var table string
	switch typ {
	case pageEntity.PageSectionItemTypeCategory:
		table = "categories"
	case pageEntity.PageSectionItemTypePost:
		table = "posts"
	case pageEntity.PageSectionItemTypeMedia:
		table = "media"
	default:
		return out, fmt.Errorf("invalid source type")
	}
	var rows []map[string]interface{}
	if err := r.db.WithContext(ctx).Table(table).Where("id IN ? AND deleted_at IS NULL", ids).Find(&rows).Error; err != nil {
		return out, err
	}
	for _, row := range rows {
		v := SourceData{}
		switch id := row["id"].(type) {
		case uint64:
			v.ID = uint(id)
		case int64:
			v.ID = uint(id)
		case uint:
			v.ID = id
		}
		fields := []struct {
			key string
			dst *string
		}{{"name", &v.Name}, {"title", &v.Title}, {"slug", &v.Slug}, {"type_code", &v.TypeCode}, {"image_url", &v.ImageURL}, {"status", &v.Status}, {"file_name", &v.FileName}, {"original_url", &v.OriginalURL}, {"thumbnail_url", &v.ThumbnailURL}, {"medium_url", &v.MediumURL}, {"mime_type", &v.MimeType}}
		for _, field := range fields {
			if value, ok := row[field.key].(string); ok {
				*field.dst = value
			}
		}
		out[v.ID] = v
	}
	return out, nil
}
