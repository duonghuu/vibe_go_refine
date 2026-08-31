package entity

import (
	"time"

	"gorm.io/gorm"
)

type PageStatus string

const (
	PageStatusDraft     PageStatus = "DRAFT"
	PageStatusPublished PageStatus = "PUBLISHED"
)

type Page struct {
	ID       uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title    string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Slug     string     `gorm:"column:slug;type:varchar(255);not null;uniqueIndex:uq_pages_slug" json:"slug"`
	Content  string     `gorm:"column:content;type:longtext;not null" json:"content"`
	Status   PageStatus `gorm:"column:status;type:varchar(20);not null;default:'DRAFT';index:idx_pages_status_deleted" json:"status"`
	AuthorID uint       `gorm:"column:author_id;not null;index:idx_pages_author_id" json:"authorId"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_pages_status_deleted" json:"-"`
}

func (Page) TableName() string {
	return "pages"
}
