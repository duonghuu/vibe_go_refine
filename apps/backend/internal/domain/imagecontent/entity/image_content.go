package entity

import (
	"time"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"gorm.io/gorm"
)

// ImageContent stores the display metadata for one image attached to an
// ImageContentType.
type ImageContent struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TypeCode string `gorm:"column:type_code;type:varchar(50);not null;index:idx_image_contents_type_status_sort,priority:1" json:"typeCode"`
	MediaID  uint   `gorm:"column:media_id;not null;index:idx_image_contents_media_id" json:"mediaId"`

	Name                 *string `gorm:"column:name;type:varchar(255)" json:"name"`
	Description          *string `gorm:"column:description;type:text" json:"description"`
	SecondaryDescription *string `gorm:"column:secondary_description;type:text" json:"secondaryDescription"`
	TargetURL            *string `gorm:"column:target_url;type:varchar(500)" json:"url"`

	SortOrder int                `gorm:"column:sort_order;not null;default:0;index:idx_image_contents_type_status_sort,priority:3" json:"sortOrder"`
	Status    ImageContentStatus `gorm:"column:status;type:varchar(20);not null;default:'ACTIVE';index:idx_image_contents_type_status_sort,priority:2" json:"status"`
	CreatedBy uint               `gorm:"column:created_by;type:int;not null;index:idx_image_contents_created_by" json:"createdBy"`

	Type      *ImageContentType  `gorm:"foreignKey:TypeCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Media     *mediaEntity.Media `gorm:"foreignKey:MediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	CreatedAt time.Time          `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt     `gorm:"column:deleted_at;index:idx_image_contents_type_status_sort,priority:4;index:idx_image_contents_deleted_at" json:"-"`
}

func (ImageContent) TableName() string { return "image_contents" }
