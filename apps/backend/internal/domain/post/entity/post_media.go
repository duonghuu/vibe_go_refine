package entity

import (
	"time"

	mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"
	"gorm.io/gorm"
)

// PostMediaCollection giới hạn các mục đích sử dụng Media trong Post.
// Alias string được dùng để giữ tương thích với các DTO hiện hữu.
type PostMediaCollection = string

const (
	PostMediaCollectionThumbnail  PostMediaCollection = "thumbnail"
	PostMediaCollectionGallery    PostMediaCollection = "gallery"
	PostMediaCollectionAttachment PostMediaCollection = "attachment"
	PostMediaCollectionContent    PostMediaCollection = "content"
)

// PostMedia đại diện cho bảng liên kết giữa Post và Media.
type PostMedia struct {
	ID         uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostID     uint           `gorm:"column:post_id;not null;index:idx_post_media_post_collection,priority:1" json:"postId"`
	MediaID    uint           `gorm:"column:media_id;not null;index:idx_post_media_media_id" json:"mediaId"`
	Collection string         `gorm:"column:collection;type:varchar(50);not null;index:idx_post_media_post_collection,priority:2" json:"collection"`
	SortOrder  int            `gorm:"column:sort_order;not null;default:0;index:idx_post_media_post_collection,priority:3" json:"sortOrder"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Post  *Post              `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
	Media *mediaEntity.Media `gorm:"foreignKey:MediaID;references:ID;constraint:OnDelete:CASCADE" json:"media,omitempty"`
}

// TableName overrides the table name used by PostMedia to `post_media`
func (PostMedia) TableName() string {
	return "post_media"
}
