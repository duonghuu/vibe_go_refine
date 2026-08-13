package entity

import (
	"time"
	"gorm.io/gorm"
)

// PostType đại diện cho bảng loại bài viết trong cơ sở dữ liệu
type PostType struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Status    string         `gorm:"type:varchar(20);default:'ACTIVE';index" json:"status"`
	SortOrder int            `gorm:"type:int;default:0;index" json:"sort_order"`

	// Audit logs & Soft Delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostType) TableName() string {
	return "post_types"
}
