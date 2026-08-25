package dto

import "time"

type CreatePostRequest struct {
	TypeCode   string `json:"typeCode" binding:"required,max=50"`
	Title      string `json:"title" binding:"required,max=255"`
	Slug       string `json:"slug" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	CategoryID *uint  `json:"categoryId"`
}

type UpdatePostRequest struct {
	Title      string `json:"title" binding:"required,max=255"`
	Slug       string `json:"slug" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	CategoryID *uint  `json:"categoryId"`
}

type GetPostByIDURI struct {
	ID uint `uri:"post_id" binding:"required,min=1"`
}

type UpdatePostURI struct {
	ID uint `uri:"post_id" binding:"required,min=1"`
}

type PostResponse struct {
	ID         uint      `json:"id"`
	TypeCode   string    `json:"typeCode"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	AuthorID   uint      `json:"authorId"`
	CategoryID *uint     `json:"categoryId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type PaginatedPostResponse struct {
	Data  []PostResponse `json:"data"`
	Total int64          `json:"total"`
}
