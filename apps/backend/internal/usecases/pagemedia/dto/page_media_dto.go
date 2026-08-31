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

type PageMediaResponse struct {
	ID         uint           `json:"id"`
	PageID     uint           `json:"pageId"`
	MediaID    uint           `json:"mediaId"`
	Collection string         `json:"collection"`
	SortOrder  int            `json:"sortOrder"`
	Media      *MediaResponse `json:"media,omitempty"`
}

type CreatePageMediaRequest struct {
	MediaID    uint   `json:"media_id" binding:"required,min=1"`
	Collection string `json:"collection" binding:"required,oneof=thumbnail gallery"`
	SortOrder  *int   `json:"sort_order" binding:"omitempty,gte=0"`
}

type MediaSyncItem struct {
	ID        uint `json:"id" binding:"required,min=1"`
	SortOrder int  `json:"sort_order" binding:"gte=0"`
}

type SyncPageMediaRequest struct {
	Media []MediaSyncItem `json:"media" binding:"max=100,dive"`
}

type GetPageMediaURI struct {
	PageID uint `uri:"page_id" binding:"required,min=1"`
}

type GetPageMediaQuery struct {
	Collection string `form:"collection" binding:"omitempty,oneof=thumbnail gallery"`
}

type SyncPageMediaURI struct {
	PageID     uint   `uri:"page_id" binding:"required,min=1"`
	Collection string `uri:"collection" binding:"required,oneof=thumbnail gallery"`
}

type DeletePageMediaURI struct {
	PageID  uint `uri:"page_id" binding:"required,min=1"`
	MediaID uint `uri:"media_id" binding:"required,min=1"`
}

type DeletePageMediaQuery struct {
	Collection string `form:"collection" binding:"omitempty,oneof=thumbnail gallery"`
}

func MapMediaResponse(media *mediaEntity.Media) *MediaResponse {
	if media == nil {
		return nil
	}
	return &MediaResponse{ID: media.ID, FileName: media.FileName, OriginalURL: media.OriginalUrl, ThumbnailURL: media.ThumbnailUrl, MediumURL: media.MediumUrl, MimeType: media.MimeType, Size: media.Size, Status: string(media.Status)}
}
