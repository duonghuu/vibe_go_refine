package dto

import "time"

type UploadMediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalUrl  string `json:"originalUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	MediumUrl    string `json:"mediumUrl"`
	Status       string `json:"status"`
}

type GetMediaListQuery struct {
	Current  int    `form:"current" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Q        string `form:"q" binding:"omitempty,max=255"`
	MimeType string `form:"mimeType" binding:"omitempty,oneof=image"`
	Status   string `form:"status" binding:"omitempty,oneof=temporary attached"`
}
type MediaListItemResponse struct {
	ID           uint      `json:"id"`
	FileName     string    `json:"fileName"`
	OriginalURL  string    `json:"originalUrl"`
	ThumbnailURL string    `json:"thumbnailUrl,omitempty"`
	MediumURL    string    `json:"mediumUrl,omitempty"`
	MimeType     string    `json:"mimeType"`
	Size         int64     `json:"size"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}
type MediaListResponse struct {
	Data  []MediaListItemResponse `json:"data"`
	Total int64                   `json:"total"`
}
