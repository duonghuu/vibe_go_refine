package repository

import (
	"context"
	"encoding/json"
	"strings"

	imageContentEntity "go_refine_dashboard_be/internal/domain/imagecontent/entity"
	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ImageContentOrderUpdate struct {
	ID        uint
	SortOrder int
}

type TypeListFilter struct {
	Offset int
	Limit  int
	Search string
	Status string
	SortBy string
	Order  string
}

type ContentListFilter struct {
	Offset   int
	Limit    int
	Search   string
	TypeCode string
	Status   string
	SortBy   string
	Order    string
}

type PublicImageContentTypeRow struct {
	Type  imageContentEntity.ImageContentType
	Items []imageContentEntity.ImageContent
}

type ImageContentRepository interface {
	FindTypeForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*imageContentEntity.ImageContentType, error)
	FindTypeByCodeForUpdate(ctx context.Context, tx *gorm.DB, code string) (*imageContentEntity.ImageContentType, error)
	FindContentTypeCode(ctx context.Context, tx *gorm.DB, id uint) (string, error)
	FindContentForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*imageContentEntity.ImageContent, error)
	FindContentsByTypeForUpdate(ctx context.Context, tx *gorm.DB, typeCode string) ([]imageContentEntity.ImageContent, error)
	FindMediaForUpdate(ctx context.Context, tx *gorm.DB, mediaID uint) (*mediaEntity.Media, error)

	CreateOrRestoreType(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContentType) error
	UpdateType(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContentType) error
	SoftDeleteType(ctx context.Context, tx *gorm.DB, id uint) error
	CreateContent(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContent) error
	UpdateContent(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContent) error
	SoftDeleteContent(ctx context.Context, tx *gorm.DB, id uint) error
	ShiftContentOrders(ctx context.Context, tx *gorm.DB, typeCode string, fromSortOrder int) error
	NormalizeContentOrders(ctx context.Context, tx *gorm.DB, items []imageContentEntity.ImageContent) error
	BulkUpdateContentOrder(ctx context.Context, tx *gorm.DB, typeCode string, updates []ImageContentOrderUpdate) error
	MarkMediaAttached(ctx context.Context, tx *gorm.DB, mediaIDs []uint) error

	FindTypes(ctx context.Context, filter TypeListFilter) ([]imageContentEntity.ImageContentType, map[string]int64, int64, error)
	FindTypeByID(ctx context.Context, id uint) (*imageContentEntity.ImageContentType, int64, error)
	FindContents(ctx context.Context, filter ContentListFilter) ([]imageContentEntity.ImageContent, int64, error)
	FindContentByID(ctx context.Context, id uint) (*imageContentEntity.ImageContent, error)
	FindPublicProjection(ctx context.Context) ([]PublicImageContentTypeRow, error)
}

type imageContentRepository struct {
	db *gorm.DB
}

func NewImageContentRepository(db *gorm.DB) ImageContentRepository {
	return &imageContentRepository{db: db}
}

func (r *imageContentRepository) FindTypeForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*imageContentEntity.ImageContentType, error) {
	var value imageContentEntity.ImageContentType
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&value, id).Error
	return &value, err
}

// FindTypeByCodeForUpdate intentionally includes soft-deleted rows so the
// command service can restore a deleted code instead of violating the unique
// constraint.
func (r *imageContentRepository) FindTypeByCodeForUpdate(ctx context.Context, tx *gorm.DB, code string) (*imageContentEntity.ImageContentType, error) {
	var value imageContentEntity.ImageContentType
	err := tx.WithContext(ctx).
		Unscoped().
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("code = ?", code).
		First(&value).Error
	return &value, err
}

func (r *imageContentRepository) FindContentTypeCode(ctx context.Context, tx *gorm.DB, id uint) (string, error) {
	var value imageContentEntity.ImageContent
	err := tx.WithContext(ctx).Select("type_code").First(&value, id).Error
	return value.TypeCode, err
}

func (r *imageContentRepository) FindContentForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*imageContentEntity.ImageContent, error) {
	var value imageContentEntity.ImageContent
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&value, id).Error
	return &value, err
}

func (r *imageContentRepository) FindContentsByTypeForUpdate(ctx context.Context, tx *gorm.DB, typeCode string) ([]imageContentEntity.ImageContent, error) {
	var values []imageContentEntity.ImageContent
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("type_code = ?", typeCode).
		Order("sort_order ASC, id ASC").
		Find(&values).Error
	return values, err
}

func (r *imageContentRepository) FindMediaForUpdate(ctx context.Context, tx *gorm.DB, mediaID uint) (*mediaEntity.Media, error) {
	var value mediaEntity.Media
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&value, mediaID).Error
	return &value, err
}

func (r *imageContentRepository) CreateOrRestoreType(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContentType) error {
	fieldConfig, err := json.Marshal(value.FieldConfig)
	if err != nil {
		return err
	}
	if value.DeletedAt.Valid {
		return tx.WithContext(ctx).Unscoped().
			Model(&imageContentEntity.ImageContentType{}).
			Where("id = ?", value.ID).
			Updates(map[string]interface{}{
				"code":         value.Code,
				"name":         value.Name,
				"status":       value.Status,
				"sort_order":   value.SortOrder,
				"max_items":    value.MaxItems,
				"field_config": fieldConfig,
				"deleted_at":   nil,
			}).Error
	}
	return tx.WithContext(ctx).Create(value).Error
}

func (r *imageContentRepository) UpdateType(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContentType) error {
	fieldConfig, err := json.Marshal(value.FieldConfig)
	if err != nil {
		return err
	}
	return tx.WithContext(ctx).Model(&imageContentEntity.ImageContentType{}).
		Where("id = ?", value.ID).
		Updates(map[string]interface{}{
			"name":         value.Name,
			"status":       value.Status,
			"sort_order":   value.SortOrder,
			"max_items":    value.MaxItems,
			"field_config": fieldConfig,
		}).Error
}

func (r *imageContentRepository) SoftDeleteType(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&imageContentEntity.ImageContentType{}, id).Error
}

func (r *imageContentRepository) CreateContent(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContent) error {
	return tx.WithContext(ctx).Create(value).Error
}

func (r *imageContentRepository) UpdateContent(ctx context.Context, tx *gorm.DB, value *imageContentEntity.ImageContent) error {
	return tx.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
		Where("id = ?", value.ID).
		Updates(map[string]interface{}{
			"media_id":              value.MediaID,
			"name":                  value.Name,
			"description":           value.Description,
			"secondary_description": value.SecondaryDescription,
			"target_url":            value.TargetURL,
			"status":                value.Status,
		}).Error
}

func (r *imageContentRepository) SoftDeleteContent(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&imageContentEntity.ImageContent{}, id).Error
}

func (r *imageContentRepository) ShiftContentOrders(ctx context.Context, tx *gorm.DB, typeCode string, fromSortOrder int) error {
	return tx.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
		Where("type_code = ? AND sort_order >= ?", typeCode, fromSortOrder).
		UpdateColumn("sort_order", gorm.Expr("sort_order + 1")).Error
}

func (r *imageContentRepository) NormalizeContentOrders(ctx context.Context, tx *gorm.DB, items []imageContentEntity.ImageContent) error {
	for index, item := range items {
		if err := tx.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
			Where("id = ?", item.ID).
			Update("sort_order", index).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *imageContentRepository) BulkUpdateContentOrder(ctx context.Context, tx *gorm.DB, typeCode string, updates []ImageContentOrderUpdate) error {
	for _, update := range updates {
		if err := tx.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
			Where("id = ? AND type_code = ?", update.ID, typeCode).
			Update("sort_order", update.SortOrder).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *imageContentRepository) MarkMediaAttached(ctx context.Context, tx *gorm.DB, mediaIDs []uint) error {
	if len(mediaIDs) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Model(&mediaEntity.Media{}).
		Where("id IN ?", mediaIDs).
		Update("status", mediaEntity.MediaStatusAttached).Error
}

func (r *imageContentRepository) FindTypes(ctx context.Context, filter TypeListFilter) ([]imageContentEntity.ImageContentType, map[string]int64, int64, error) {
	query := r.db.WithContext(ctx).Model(&imageContentEntity.ImageContentType{})
	if filter.Search != "" {
		pattern := "%" + escapeLike(filter.Search) + "%"
		query = query.Where("(code LIKE ? OR name LIKE ?)", pattern, pattern)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	counts := make(map[string]int64)
	var grouped []struct {
		TypeCode string
		Count    int64
	}
	if err := r.db.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
		Select("type_code, COUNT(*) AS count").
		Group("type_code").
		Scan(&grouped).Error; err != nil {
		return nil, nil, 0, err
	}
	for _, item := range grouped {
		counts[item.TypeCode] = item.Count
	}

	sortColumn := typeSortColumn(filter.SortBy)
	order := normalizeOrder(filter.Order)
	var values []imageContentEntity.ImageContentType
	if err := query.Order(sortColumn + " " + order + ", id ASC").
		Offset(filter.Offset).Limit(filter.Limit).Find(&values).Error; err != nil {
		return nil, nil, 0, err
	}
	return values, counts, total, nil
}

func (r *imageContentRepository) FindTypeByID(ctx context.Context, id uint) (*imageContentEntity.ImageContentType, int64, error) {
	var value imageContentEntity.ImageContentType
	if err := r.db.WithContext(ctx).First(&value, id).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
		Where("type_code = ?", value.Code).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return &value, count, nil
}

func (r *imageContentRepository) FindContents(ctx context.Context, filter ContentListFilter) ([]imageContentEntity.ImageContent, int64, error) {
	query := r.db.WithContext(ctx).Model(&imageContentEntity.ImageContent{}).
		Preload("Type").Preload("Media")
	if filter.Search != "" {
		pattern := "%" + escapeLike(filter.Search) + "%"
		query = query.Where("(name LIKE ? OR description LIKE ? OR secondary_description LIKE ? OR target_url LIKE ?)", pattern, pattern, pattern, pattern)
	}
	if filter.TypeCode != "" {
		query = query.Where("type_code = ?", filter.TypeCode)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	sortColumn := contentSortColumn(filter.SortBy)
	order := normalizeOrder(filter.Order)
	var values []imageContentEntity.ImageContent
	if err := query.Order(sortColumn + " " + order + ", id ASC").
		Offset(filter.Offset).Limit(filter.Limit).Find(&values).Error; err != nil {
		return nil, 0, err
	}
	return values, total, nil
}

func (r *imageContentRepository) FindContentByID(ctx context.Context, id uint) (*imageContentEntity.ImageContent, error) {
	var value imageContentEntity.ImageContent
	err := r.db.WithContext(ctx).Preload("Type").Preload("Media").First(&value, id).Error
	return &value, err
}

func (r *imageContentRepository) FindPublicProjection(ctx context.Context) ([]PublicImageContentTypeRow, error) {
	var types []imageContentEntity.ImageContentType
	if err := r.db.WithContext(ctx).
		Where("status = ?", imageContentEntity.ImageContentStatusActive).
		Order("sort_order ASC, id ASC").
		Find(&types).Error; err != nil {
		return nil, err
	}

	rows := make([]PublicImageContentTypeRow, 0, len(types))
	if len(types) == 0 {
		return rows, nil
	}
	var items []imageContentEntity.ImageContent
	if err := r.db.WithContext(ctx).
		Preload("Media").
		Joins("JOIN image_content_types ON image_content_types.code = image_contents.type_code AND image_content_types.deleted_at IS NULL").
		Joins("JOIN media ON media.id = image_contents.media_id AND media.deleted_at IS NULL").
		Where("image_content_types.status = ?", imageContentEntity.ImageContentStatusActive).
		Where("image_contents.status = ?", imageContentEntity.ImageContentStatusActive).
		Where("media.status = ?", mediaEntity.MediaStatusAttached).
		Where("media.mime_type LIKE ?", "image/%").
		Order("image_contents.type_code ASC, image_contents.sort_order ASC, image_contents.id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	byCode := make(map[string][]imageContentEntity.ImageContent, len(types))
	for _, item := range items {
		byCode[item.TypeCode] = append(byCode[item.TypeCode], item)
	}
	for _, value := range types {
		rows = append(rows, PublicImageContentTypeRow{Type: value, Items: byCode[value.Code]})
	}
	return rows, nil
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	return strings.ReplaceAll(value, `_`, `\_`)
}

func normalizeOrder(value string) string {
	if strings.EqualFold(value, "DESC") {
		return "DESC"
	}
	return "ASC"
}

func typeSortColumn(value string) string {
	switch value {
	case "id":
		return "id"
	case "code":
		return "code"
	case "name":
		return "name"
	case "status":
		return "status"
	case "createdAt":
		return "created_at"
	case "updatedAt":
		return "updated_at"
	default:
		return "sort_order"
	}
}

func contentSortColumn(value string) string {
	switch value {
	case "id":
		return "id"
	case "typeCode":
		return "type_code"
	case "name":
		return "name"
	case "status":
		return "status"
	case "createdAt":
		return "created_at"
	case "updatedAt":
		return "updated_at"
	default:
		return "sort_order"
	}
}
