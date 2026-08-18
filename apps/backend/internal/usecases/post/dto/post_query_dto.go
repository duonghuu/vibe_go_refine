package dto

type GetPostsQuery struct {
	Page     int    `form:"current" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1"`
	Title    string `form:"title_like" binding:"omitempty"`
	TypeCode string `form:"typeCode" binding:"omitempty,max=50"`
	SortBy   string `form:"sortBy" binding:"omitempty"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}
