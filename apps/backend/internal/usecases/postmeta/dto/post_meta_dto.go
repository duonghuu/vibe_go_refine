package dto

import "time"

type MetaItemRequest struct {
	Key   string `json:"key" binding:"required,max=255"`
	Value string `json:"value"`
}

type SyncPostMetaRequest struct {
	PostID uint              `uri:"postId" binding:"required"`
	Meta   []MetaItemRequest `json:"meta" binding:"required,dive"`
}

type GetPostMetaRequest struct {
	PostID uint `uri:"postId" binding:"required"`
}

type DeletePostMetaRequest struct {
	PostID uint   `uri:"postId" binding:"required"`
	Key    string `uri:"key" binding:"required"`
}

type PostMetaResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetPostMetaResponse struct {
	Data []PostMetaResponse `json:"data"`
}

type SyncPostMetaResponse struct {
	Message string `json:"message"`
}

type DeletePostMetaResponse struct {
	Message string `json:"message"`
}
