package dto

type PostMediaResponse struct {
	ID         uint        `json:"id"`
	PostID     uint        `json:"post_id"`
	MediaID    uint        `json:"media_id"`
	Collection string      `json:"collection"`
	SortOrder  int         `json:"sort_order"`
	Media      interface{} `json:"media,omitempty"`
}

type CreatePostMediaRequest struct {
	MediaID    uint   `json:"media_id" binding:"required"`
	Collection string `json:"collection" binding:"required,oneof=thumbnail gallery attachment content"`
	SortOrder  *int   `json:"sort_order"`
}

type SyncPostMediaRequest struct {
	Media []MediaSyncItem `json:"media" binding:"required,dive"`
}

type MediaSyncItem struct {
	ID        uint `json:"id" binding:"required"`
	SortOrder int  `json:"sort_order"`
}
