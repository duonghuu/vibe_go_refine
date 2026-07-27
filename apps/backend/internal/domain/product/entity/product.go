package entity

import (
	"time"

	categoryEntity "go_refine_dashboard_be/internal/domain/category/entity"

	"gorm.io/gorm"
)

// Product đại diện cho thực thể Sản phẩm
type Product struct {
	ID         uint                    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	SKU        string                  `gorm:"type:varchar(50);uniqueIndex;not null;column:sku" json:"sku"`
	Slug       string                  `gorm:"type:varchar(255);uniqueIndex;not null;column:slug" json:"slug"`
	Name       string                  `gorm:"type:varchar(255);not null;index:idx_product_name;column:name" json:"name"`
	Description string                 `gorm:"type:text;column:description" json:"description"`
	CategoryID uint                    `gorm:"not null;index;column:category_id" json:"categoryId"`
	Category   categoryEntity.Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:RESTRICT;" json:"category"`
	Price      float64                 `gorm:"type:decimal(10,2);not null;column:price" json:"price"`
	SalePrice  *float64                `gorm:"type:decimal(10,2);column:sale_price" json:"salePrice"`
	Stock      int                     `gorm:"type:int;not null;default:0;column:stock" json:"stock"`
	SoldCount  int                     `gorm:"type:int;not null;default:0;column:sold_count" json:"soldCount"`
	Status     string                  `gorm:"type:varchar(20);not null;default:'ACTIVE';index;column:status" json:"status"`
	ImageURL   string                  `gorm:"type:text;column:image_url" json:"image"`
	CreatedAt  time.Time               `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdatedAt  time.Time               `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt          `gorm:"index;column:deleted_at" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
