package entity

import (
	"time"
	"gorm.io/gorm"
)

// Category đại diện cho bảng danh mục sản phẩm trong cơ sở dữ liệu
type Category struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	ParentID    *uint          `gorm:"index" json:"parentId"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"imageUrl"`
	SortOrder   int            `gorm:"type:int;default:0;index" json:"sortOrder"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"` // ACTIVE, HIDDEN
	
	// Quan hệ nội bộ (Self-referencing)
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Audit logs & Soft Delete
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}
