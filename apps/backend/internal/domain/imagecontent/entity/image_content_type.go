package entity

import (
	"time"

	"gorm.io/gorm"
)

// ImageContentType defines the metadata configuration and item limit for a
// site-wide image content collection.
type ImageContentType struct {
	ID          uint                    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string                  `gorm:"column:code;type:varchar(50);not null;uniqueIndex:uq_image_content_types_code" json:"code"`
	Name        string                  `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Status      ImageContentStatus      `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE';index:idx_image_content_types_status_sort,priority:1" json:"status"`
	SortOrder   int                     `gorm:"column:sort_order;not null;default:0;index:idx_image_content_types_status_sort,priority:2" json:"sortOrder"`
	MaxItems    *uint                   `gorm:"column:max_items;type:int unsigned" json:"maxItems"`
	FieldConfig ImageContentFieldConfig `gorm:"column:field_config;type:json;serializer:json;not null" json:"-"`

	Items     []ImageContent `gorm:"foreignKey:TypeCode;references:Code" json:"-"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_image_content_types_deleted_at" json:"-"`
}

func (ImageContentType) TableName() string { return "image_content_types" }
