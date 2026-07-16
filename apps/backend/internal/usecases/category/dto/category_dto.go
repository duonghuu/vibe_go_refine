package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Slug        *string `json:"slug" binding:"omitempty,max=255"`
	ParentID    *uint   `json:"parentId"`
	Description *string `json:"description"`
	ImageURL    *string `json:"imageUrl"`
	SortOrder   *int    `json:"sortOrder"`
	Status      *string `json:"status" binding:"omitempty,oneof=ACTIVE HIDDEN"`
}

type BulkUpdateCategoryStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status string `json:"status" binding:"required,oneof=ACTIVE HIDDEN"`
}

type ReorderCategoryItem struct {
	ID        uint `json:"id" binding:"required"`
	SortOrder int  `json:"sortOrder" binding:"required"`
}

type ReorderCategoriesRequest struct {
	Items []ReorderCategoryItem `json:"items" binding:"required,min=1"`
}

type CategoryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     *uint  `json:"parentId"`
	Description  string `json:"description"`
	ImageURL     string `json:"imageUrl"`
	SortOrder    int    `json:"sortOrder"`
	Status       string `json:"status"`
	ProductCount int64  `json:"productCount"`
	CreatedAt    string `json:"createdAt"`
}
