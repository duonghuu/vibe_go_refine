package entity

import (
	"time"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"gorm.io/gorm"
)

type PageMediaCollection string

const (
	PageMediaCollectionThumbnail PageMediaCollection = "thumbnail"
	PageMediaCollectionGallery   PageMediaCollection = "gallery"
)

// PageMedia represents the link between a Page and a Media asset.
// Soft deletion applies to the link only; the original Media asset is preserved.
type PageMedia struct {
	ID         uint                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PageID     uint                `gorm:"column:page_id;not null;index:idx_page_media_page_collection_sort,priority:1;uniqueIndex:uq_page_media_link,priority:1" json:"pageId"`
	MediaID    uint                `gorm:"column:media_id;not null;index:idx_page_media_media_id;uniqueIndex:uq_page_media_link,priority:2" json:"mediaId"`
	Collection PageMediaCollection `gorm:"column:collection;type:varchar(20);not null;index:idx_page_media_page_collection_sort,priority:2;uniqueIndex:uq_page_media_link,priority:3" json:"collection"`
	SortOrder  int                 `gorm:"column:sort_order;not null;default:0;index:idx_page_media_page_collection_sort,priority:3" json:"sortOrder"`

	Media     *mediaEntity.Media `gorm:"foreignKey:MediaID;references:ID;constraint:OnDelete:CASCADE" json:"media,omitempty"`
	CreatedAt time.Time          `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt     `gorm:"column:deleted_at;index:idx_page_media_page_collection_sort,priority:4" json:"-"`
}

func (PageMedia) TableName() string {
	return "page_media"
}
