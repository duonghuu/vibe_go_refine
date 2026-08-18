package dto

import "time"

type CreatePostRequest struct {
	TypeCode string `json:"typeCode" binding:"required,max=50"`
	Title    string `json:"title" binding:"required,max=255"`
	Slug     string `json:"slug" binding:"required,max=255"`
	Content  string `json:"content" binding:"required"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	TypeCode  string    `json:"typeCode"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	AuthorID  uint      `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedPostResponse struct {
	Data  []PostResponse `json:"data"`
	Total int64          `json:"total"`
}
