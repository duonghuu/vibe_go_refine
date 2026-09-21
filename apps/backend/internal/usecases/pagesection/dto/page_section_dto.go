package dto

import "time"

type PageSectionsURI struct {
	PageID uint `uri:"page_id" binding:"required,min=1"`
}
type PageSectionURI struct {
	PageID    uint `uri:"page_id" binding:"required,min=1"`
	SectionID uint `uri:"section_id" binding:"required,min=1"`
}
type SectionCollectionURI struct {
	PageID     uint   `uri:"page_id" binding:"required,min=1"`
	SectionID  uint   `uri:"section_id" binding:"required,min=1"`
	Collection string `uri:"collection" binding:"required,max=100"`
}
type GetPageSectionsQuery struct {
	Current  int    `form:"current" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Status   string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
type GetSectionItemsQuery struct {
	Collection string `form:"collection" binding:"required,max=100"`
	Current    int    `form:"current" binding:"omitempty,min=1"`
	PageSize   int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
}
type CreatePageSectionRequest struct {
	Key               string  `json:"key" binding:"required,max=100"`
	Name              string  `json:"name" binding:"required,max=255"`
	Title             string  `json:"title" binding:"omitempty,max=255"`
	Description       string  `json:"description" binding:"omitempty,max=5000"`
	BackgroundColor   *string `json:"backgroundColor" binding:"omitempty,max=9"`
	BackgroundMediaID *uint   `json:"backgroundMediaId" binding:"omitempty,min=1"`
	FeatureMediaID    *uint   `json:"featureMediaId" binding:"omitempty,min=1"`
	SortOrder         *int    `json:"sortOrder" binding:"omitempty,gte=0"`
	Status            string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}
type UpdatePageSectionRequest struct {
	Key               string  `json:"key" binding:"required,max=100"`
	Name              string  `json:"name" binding:"required,max=255"`
	Title             string  `json:"title" binding:"omitempty,max=255"`
	Description       string  `json:"description" binding:"omitempty,max=5000"`
	BackgroundColor   *string `json:"backgroundColor" binding:"omitempty,max=9"`
	BackgroundMediaID *uint   `json:"backgroundMediaId" binding:"omitempty,min=1"`
	FeatureMediaID    *uint   `json:"featureMediaId" binding:"omitempty,min=1"`
	Status            string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}
type ReorderSectionItem struct {
	ID        uint `json:"id" binding:"required,min=1"`
	SortOrder int  `json:"sortOrder" binding:"gte=0"`
}
type ReorderPageSectionsRequest struct {
	Sections []ReorderSectionItem `json:"sections" binding:"required,min=1,max=100,dive"`
}
type SyncSectionItem struct {
	ItemID    uint `json:"itemId" binding:"required,min=1"`
	SortOrder int  `json:"sortOrder" binding:"gte=0"`
}
type SyncSectionCollectionRequest struct {
	ItemType string            `json:"itemType" binding:"required,oneof=CATEGORY POST MEDIA"`
	Items    []SyncSectionItem `json:"items" binding:"max=100,dive"`
}
type SectionMediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
	MimeType     string `json:"mimeType"`
}
type SectionCollectionSummary struct {
	Collection string `json:"collection"`
	ItemType   string `json:"itemType"`
	Total      int64  `json:"total"`
}
type PageSectionResponse struct {
	ID                uint                       `json:"id"`
	PageID            uint                       `json:"pageId"`
	Key               string                     `json:"key"`
	Name              string                     `json:"name"`
	Title             string                     `json:"title"`
	Description       string                     `json:"description"`
	BackgroundColor   *string                    `json:"backgroundColor"`
	BackgroundMediaID *uint                      `json:"backgroundMediaId"`
	BackgroundMedia   *SectionMediaResponse      `json:"backgroundMedia"`
	FeatureMediaID    *uint                      `json:"featureMediaId"`
	FeatureMedia      *SectionMediaResponse      `json:"featureMedia"`
	SortOrder         int                        `json:"sortOrder"`
	Status            string                     `json:"status"`
	Collections       []SectionCollectionSummary `json:"collections"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
}
type PageSectionListResponse struct {
	Data  []PageSectionResponse `json:"data"`
	Total int64                 `json:"total"`
}
type SectionItemDataResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description"`
	Slug         string `json:"slug,omitempty"`
	TypeCode     string `json:"typeCode,omitempty"`
	ImageURL     string `json:"imageUrl,omitempty"`
	Status       string `json:"status,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	OriginalURL  string `json:"originalUrl,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
	MimeType     string `json:"mimeType,omitempty"`
}
type PageSectionItemResponse struct {
	ID         uint                    `json:"id"`
	SectionID  uint                    `json:"sectionId"`
	ItemType   string                  `json:"itemType"`
	ItemID     uint                    `json:"itemId"`
	Collection string                  `json:"collection"`
	SortOrder  int                     `json:"sortOrder"`
	Data       SectionItemDataResponse `json:"data"`
}
type PageSectionItemListResponse struct {
	Data  []PageSectionItemResponse `json:"data"`
	Total int64                     `json:"total"`
}
