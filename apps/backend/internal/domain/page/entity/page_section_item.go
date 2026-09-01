package entity

import (
	"time"

	"gorm.io/gorm"
)

type PageSectionItemType string

const (
	PageSectionItemTypeCategory PageSectionItemType = "CATEGORY"
	PageSectionItemTypePost     PageSectionItemType = "POST"
	PageSectionItemTypeMedia    PageSectionItemType = "MEDIA"
)

type PageSectionItem struct {
	ID         uint                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SectionID  uint                `gorm:"column:section_id;not null;uniqueIndex:uq_page_section_item_link,priority:1;index:idx_page_section_items_section_collection_sort,priority:1" json:"sectionId"`
	ItemType   PageSectionItemType `gorm:"column:item_type;type:varchar(20);not null;uniqueIndex:uq_page_section_item_link,priority:3;index:idx_page_section_items_source,priority:1" json:"itemType"`
	ItemID     uint                `gorm:"column:item_id;not null;uniqueIndex:uq_page_section_item_link,priority:4;index:idx_page_section_items_source,priority:2" json:"itemId"`
	Collection string              `gorm:"column:collection;type:varchar(100);not null;uniqueIndex:uq_page_section_item_link,priority:2;index:idx_page_section_items_section_collection_sort,priority:2" json:"collection"`
	SortOrder  int                 `gorm:"column:sort_order;not null;default:0;index:idx_page_section_items_section_collection_sort,priority:3" json:"sortOrder"`

	Section   *PageSection   `gorm:"foreignKey:SectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_page_section_items_section_collection_sort,priority:4;index:idx_page_section_items_source,priority:3" json:"-"`
}

func (PageSectionItem) TableName() string {
	return "page_section_items"
}
