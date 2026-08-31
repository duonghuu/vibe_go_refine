package dto

import "time"

type CreatePageRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Slug    string `json:"slug" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
}

type PageResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	AuthorID  uint      `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
