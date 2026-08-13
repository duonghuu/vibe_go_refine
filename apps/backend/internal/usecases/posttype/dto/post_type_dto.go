package dto

import "time"

type PostTypeResponse struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreatePostTypeRequest struct {
	Code      string `json:"code" binding:"required,max=50,alphanum"`
	Name      string `json:"name" binding:"required,max=100"`
	Status    string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	SortOrder int    `json:"sort_order" binding:"omitempty,min=0"`
}

type UpdatePostTypeRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Status    string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	SortOrder int    `json:"sort_order" binding:"omitempty,min=0"`
}
