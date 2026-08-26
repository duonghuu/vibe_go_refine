package entity

import "time"

// SEOMeta stores the single SEO configuration owned by an entity.
// The polymorphic target is validated by the SEO service registry because
// entity_id cannot have a physical foreign key to multiple tables.
type SEOMeta struct {
	ID                 uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	EntityType         string    `gorm:"column:entity_type;type:varchar(50);not null" json:"entity_type"`
	EntityID           uint      `gorm:"column:entity_id;not null" json:"entity_id"`
	MetaTitle          *string   `gorm:"column:meta_title;type:varchar(255)" json:"meta_title"`
	MetaDescription    *string   `gorm:"column:meta_description;type:varchar(500)" json:"meta_description"`
	MetaKeywords       *string   `gorm:"column:meta_keywords;type:varchar(500)" json:"meta_keywords"`
	CanonicalURL       *string   `gorm:"column:canonical_url;type:varchar(500)" json:"canonical_url"`
	OGTitle            *string   `gorm:"column:og_title;type:varchar(255)" json:"og_title"`
	OGDescription      *string   `gorm:"column:og_description;type:varchar(500)" json:"og_description"`
	OGImage            *string   `gorm:"column:og_image;type:varchar(500)" json:"og_image"`
	OGType             *string   `gorm:"column:og_type;type:varchar(100)" json:"og_type"`
	TwitterTitle       *string   `gorm:"column:twitter_title;type:varchar(255)" json:"twitter_title"`
	TwitterDescription *string   `gorm:"column:twitter_description;type:varchar(500)" json:"twitter_description"`
	TwitterImage       *string   `gorm:"column:twitter_image;type:varchar(500)" json:"twitter_image"`
	TwitterCard        *string   `gorm:"column:twitter_card;type:varchar(100)" json:"twitter_card"`
	Robots             *string   `gorm:"column:robots;type:varchar(100)" json:"robots"`
	SchemaJSON         []byte    `gorm:"column:schema_json;type:json" json:"schema_json"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SEOMeta) TableName() string {
	return "seo_meta"
}
