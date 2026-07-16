package entity

import (
	"time"

	productEntity "go_refine_dashboard_be/internal/domain/product/entity"

	"gorm.io/gorm"
)

type MediaStatus string

const (
	MediaStatusTemporary MediaStatus = "temporary"
	MediaStatusAttached  MediaStatus = "attached"
)

// Media đại diện cho thực thể File lưu trữ trên hệ thống
type Media struct {
	ID           uint                   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	OriginalName string                 `gorm:"type:varchar(255);not null;column:original_name" json:"originalName"`
	FileName     string                 `gorm:"type:varchar(255);not null;uniqueIndex;column:file_name" json:"fileName"`
	MimeType     string                 `gorm:"type:varchar(100);not null;column:mime_type" json:"mimeType"`
	Size         int64                  `gorm:"not null;column:size" json:"size"`
	OriginalUrl  string                 `gorm:"type:varchar(500);not null;column:original_url" json:"originalUrl"`
	ThumbnailUrl string                 `gorm:"type:varchar(500);column:thumbnail_url" json:"thumbnailUrl"`
	MediumUrl    string                 `gorm:"type:varchar(500);column:medium_url" json:"mediumUrl"`
	Status       MediaStatus            `gorm:"type:enum('temporary', 'attached');default:'temporary';index:idx_status_created;column:status" json:"status"`

	OwnerID      uint                   `gorm:"not null;index;column:owner_id" json:"ownerId"`
	ProductID    *uint                  `gorm:"index;column:product_id" json:"productId"`
	Product      *productEntity.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:SET NULL;" json:"product,omitempty"`

	CreatedAt    time.Time              `gorm:"autoCreateTime;index:idx_status_created;column:created_at" json:"createdAt"`
	UpdatedAt    time.Time              `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt         `gorm:"index;column:deleted_at" json:"-"`
}

func (Media) TableName() string {
	return "media"
}
