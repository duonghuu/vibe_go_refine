package entity

import (
	"time"
	"gorm.io/gorm"
)

// PostCategory đại diện cho bảng danh mục bài viết trong cơ sở dữ liệu
type PostCategory struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode    string         `gorm:"type:varchar(50);not null;index" json:"typeCode"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	ParentID    *uint          `gorm:"index" json:"parentId"` // Nullable để hỗ trợ danh mục gốc
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"imageUrl"`
	SortOrder   int            `gorm:"type:int;default:0;index" json:"sortOrder"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, INACTIVE
	
	// Quan hệ nội bộ (Self-referencing) để tạo dạng cây
	Parent      *PostCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []PostCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Audit logs & Soft Delete
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostCategory) TableName() string {
	return "post_categories"
}
