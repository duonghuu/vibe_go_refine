package dto

import mediaEntity "go_refine_dashboard_be/internal/domain/media/entity"

type MediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
	MimeType     string `json:"mimeType"`
	Size         int64  `json:"size"`
	Status       string `json:"status"`
}

type PostMediaResponse struct {
	ID         uint           `json:"id"`
	PostID     uint           `json:"postId"`
	MediaID    uint           `json:"mediaId"`
	Collection string         `json:"collection"`
	SortOrder  int            `json:"sortOrder"`
	Media      *MediaResponse `json:"media,omitempty"`
}

type CreatePostMediaRequest struct {
	MediaID    uint   `json:"media_id" binding:"required,min=1"`
	Collection string `json:"collection" binding:"required,oneof=thumbnail gallery"`
	SortOrder  *int   `json:"sort_order" binding:"omitempty,gte=0"`
}

type SyncPostMediaRequest struct {
	Media []MediaSyncItem `json:"media" binding:"max=100,dive"`
}

type MediaSyncItem struct {
	ID        uint `json:"id" binding:"required,min=1"`
	SortOrder int  `json:"sort_order" binding:"gte=0"`
}

type MediaPayload struct {
	ThumbnailID *uint           `json:"thumbnailId"`
	Gallery     []MediaSyncItem `json:"gallery" binding:"max=100,dive"`
}

func MapMediaResponse(media *mediaEntity.Media) *MediaResponse {
	if media == nil {
		return nil
	}
	return &MediaResponse{
		ID:           media.ID,
		FileName:     media.FileName,
		OriginalURL:  media.OriginalUrl,
		ThumbnailURL: media.ThumbnailUrl,
		MediumURL:    media.MediumUrl,
		MimeType:     media.MimeType,
		Size:         media.Size,
		Status:       string(media.Status),
	}
}
