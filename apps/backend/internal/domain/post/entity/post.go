package entity

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeCode  string         `gorm:"type:varchar(50);not null;index" json:"typeCode"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content   string         `gorm:"type:longtext;not null" json:"content"`
	AuthorID  uint           `gorm:"not null;index" json:"authorId"`
	
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Post) TableName() string {
	return "posts"
}
