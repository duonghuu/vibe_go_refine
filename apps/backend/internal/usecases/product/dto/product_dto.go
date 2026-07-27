package dto

type ListProductReq struct {
	Page       int      `form:"_start" binding:"min=0"`
	PageSize   int      `form:"_end" binding:"min=1"`
	Sort       string   `form:"_sort"`
	Order      string   `form:"_order"`
	Search     string   `form:"q"`
	CategoryID *uint    `form:"categoryId"`
	Status     string   `form:"status"`
	MinPrice   *float64 `form:"minPrice"`
	MaxPrice   *float64 `form:"maxPrice"`
	MinStock   *int     `form:"minStock"`
}

type CreateProductReq struct {
	Name       string   `json:"name" binding:"required,max=255"`
	SKU        string   `json:"sku" binding:"required,max=50"`
	CategoryID uint     `json:"categoryId" binding:"required"`
	Price      float64  `json:"price" binding:"required,gt=0"`
	SalePrice  *float64 `json:"salePrice" binding:"omitempty,gte=0"`
	Stock      int      `json:"stock" binding:"required,gte=0"`
	Status     string   `json:"status" binding:"required,oneof=ACTIVE HIDDEN OUT_OF_STOCK"`
	ImageURL   string   `json:"image" binding:"omitempty,max=500"`
}

type UpdateProductReq struct {
	CreateProductReq // embed for now
}

type UpdateProductStatusReq struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE HIDDEN OUT_OF_STOCK"`
}

type BulkDeleteReq struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

type BulkStatusReq struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=ACTIVE HIDDEN OUT_OF_STOCK"`
}
