package dto

type CreatePostCategoryRequest struct {
	TypeCode    string `json:"typeCode" binding:"required,max=50"`
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	ParentID    *uint  `json:"parentId"` // Nullable
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder" binding:"omitempty,min=0"`
	Status      string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type PostCategoryResponse struct {
	ID          uint   `json:"id"`
	TypeCode    string `json:"typeCode"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ParentID    *uint  `json:"parentId"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

type PostCategoryTreeResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"` // Có thể được format thêm prefix "--" theo level
	TypeCode string `json:"typeCode"`
}
