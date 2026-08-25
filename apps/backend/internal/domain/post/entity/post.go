package entity

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID         uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TypeCode   string `gorm:"column:type_code;type:varchar(50);not null;index" json:"typeCode"`
	Title      string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Slug       string `gorm:"column:slug;type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content    string `gorm:"column:content;type:longtext;not null" json:"content"`
	AuthorID   uint   `gorm:"column:author_id;not null;index" json:"authorId"`
	CategoryID *uint  `gorm:"column:category_id;index" json:"categoryId"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Post) TableName() string {
	return "posts"
}
