package dto

import (
	"time"

	"go_refine_dashboard_be/internal/domain/page/entity"
)

type CreatePageRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Slug    string `json:"slug" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
}

type GetPagesQuery struct {
	Current  int    `form:"current" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Search   string `form:"title_like" binding:"omitempty,max=255"`
	Status   string `form:"status" binding:"omitempty,oneof=DRAFT PUBLISHED"`
	SortBy   string `form:"sortBy" binding:"omitempty,oneof=id title slug updated_at"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

type PageListItemResponse struct {
	ID        uint              `json:"id"`
	Title     string            `json:"title"`
	Slug      string            `json:"slug"`
	Status    entity.PageStatus `json:"status"`
	AuthorID  uint              `json:"authorId"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type PaginatedPageResponse struct {
	Data  []PageListItemResponse `json:"data"`
	Total int64                  `json:"total"`
}

type DeletePageURI struct {
	PageID uint `uri:"page_id" binding:"required,min=1"`
}

type DeletePageResponse struct {
	Message string `json:"message"`
}

type GetPageByIDURI struct {
	PageID uint `uri:"page_id" binding:"required,min=1"`
}

type UpdatePageURI struct {
	PageID uint `uri:"page_id" binding:"required,min=1"`
}

type UpdatePageRequest struct {
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
