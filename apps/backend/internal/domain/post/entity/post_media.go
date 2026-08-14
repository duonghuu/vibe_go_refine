package entity

import (
	"time"

	"gorm.io/gorm"
)

// PostMedia đại diện cho bảng liên kết giữa Post và Media
type PostMedia struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID     uint           `gorm:"not null;index:idx_post_media_post_collection" json:"post_id"`
	MediaID    uint           `gorm:"not null;index" json:"media_id"`
	Collection string         `gorm:"type:varchar(50);not null;index:idx_post_media_post_collection" json:"collection"` // thumbnail, gallery, attachment, content
	SortOrder  int            `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Define TableName
}

// TableName overrides the table name used by PostMedia to `post_media`
func (PostMedia) TableName() string {
	return "post_media"
}
