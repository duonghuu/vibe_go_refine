package entity

import (
	"time"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"gorm.io/gorm"
)

type PageSectionStatus string

const (
	PageSectionStatusActive   PageSectionStatus = "ACTIVE"
	PageSectionStatusInactive PageSectionStatus = "INACTIVE"
)

type PageSection struct {
	ID     uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PageID uint   `gorm:"column:page_id;not null;uniqueIndex:uq_page_sections_page_key,priority:1;index:idx_page_sections_page_status_sort,priority:1" json:"pageId"`
	Key    string `gorm:"column:key;type:varchar(100);not null;uniqueIndex:uq_page_sections_page_key,priority:2" json:"key"`
	Name   string `gorm:"column:name;type:varchar(255);not null" json:"name"`

	Title       string `gorm:"column:title;type:varchar(255);not null;default:''" json:"title"`
	Description string `gorm:"column:description;type:text;not null" json:"description"`

	BackgroundColor   *string `gorm:"column:background_color;type:varchar(9)" json:"backgroundColor"`
	BackgroundMediaID *uint   `gorm:"column:background_media_id;index:idx_page_sections_background_media" json:"backgroundMediaId"`
	FeatureMediaID    *uint   `gorm:"column:feature_media_id;index:idx_page_sections_feature_media" json:"featureMediaId"`

	SortOrder int               `gorm:"column:sort_order;not null;default:0;index:idx_page_sections_page_status_sort,priority:3" json:"sortOrder"`
	Status    PageSectionStatus `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE';index:idx_page_sections_page_status_sort,priority:2" json:"status"`

	Page            *Page              `gorm:"foreignKey:PageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	BackgroundMedia *mediaEntity.Media `gorm:"foreignKey:BackgroundMediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"backgroundMedia,omitempty"`
	FeatureMedia    *mediaEntity.Media `gorm:"foreignKey:FeatureMediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"featureMedia,omitempty"`
	Items           []PageSectionItem  `gorm:"foreignKey:SectionID;references:ID" json:"-"`
	CreatedAt       time.Time          `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt     `gorm:"column:deleted_at;index:idx_page_sections_page_status_sort,priority:4" json:"-"`
}

func (PageSection) TableName() string {
	return "page_sections"
}
