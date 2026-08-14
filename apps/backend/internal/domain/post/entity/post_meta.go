package entity

import (
	"time"
)

type PostMeta struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index:idx_post_meta_post_key,unique" json:"post_id"`
	Key       string    `gorm:"type:varchar(255);not null;index:idx_post_meta_post_key,unique" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PostMeta) TableName() string {
	return "post_meta"
}
