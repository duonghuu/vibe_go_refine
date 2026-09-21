package dto

import (
	"encoding/json"
	"time"
)

type GetPublicPageURI struct {
	Slug string `uri:"slug" binding:"required,max=255"`
}

type PublicPageMediaResponse struct {
	ID           uint   `json:"id"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
}

type PublicPageSEOResponse struct {
	ResolvedTitle        string          `json:"resolvedTitle"`
	ResolvedDescription  string          `json:"resolvedDescription"`
	ResolvedCanonicalURL string          `json:"resolvedCanonicalUrl"`
	OGTitle              string          `json:"ogTitle,omitempty"`
	OGDescription        string          `json:"ogDescription,omitempty"`
	OGImage              string          `json:"ogImage,omitempty"`
	OGType               string          `json:"ogType,omitempty"`
	TwitterTitle         string          `json:"twitterTitle,omitempty"`
	TwitterDescription   string          `json:"twitterDescription,omitempty"`
	TwitterImage         string          `json:"twitterImage,omitempty"`
	TwitterCard          string          `json:"twitterCard,omitempty"`
	Robots               string          `json:"robots"`
	SchemaJSON           json.RawMessage `json:"schemaJson,omitempty"`
}

type PublicSectionItemDataResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description"`
	Title        string `json:"title,omitempty"`
	Slug         string `json:"slug,omitempty"`
	TypeCode     string `json:"typeCode,omitempty"`
	ImageURL     string `json:"imageUrl,omitempty"`
	OriginalURL  string `json:"originalUrl,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
}

type PublicPageSectionItemResponse struct {
	ItemType  string                        `json:"itemType"`
	ItemID    uint                          `json:"itemId"`
	SortOrder int                           `json:"sortOrder"`
	Data      PublicSectionItemDataResponse `json:"data"`
}

type PublicPageSectionMediaResponse struct {
	ID           uint   `json:"id"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
}

type PublicPageSectionResponse struct {
	ID              uint                                       `json:"id"`
	Key             string                                     `json:"key"`
	Title           string                                     `json:"title"`
	Description     string                                     `json:"description"`
	BackgroundColor *string                                    `json:"backgroundColor"`
	BackgroundMedia *PublicPageSectionMediaResponse            `json:"backgroundMedia"`
	FeatureMedia    *PublicPageSectionMediaResponse            `json:"featureMedia"`
	SortOrder       int                                        `json:"sortOrder"`
	Collections     map[string][]PublicPageSectionItemResponse `json:"collections"`
}

type PublicPageResponse struct {
	ID        uint                        `json:"id"`
	Title     string                      `json:"title"`
	Slug      string                      `json:"slug"`
	Content   string                      `json:"content"`
	UpdatedAt time.Time                   `json:"updatedAt"`
	Thumbnail *PublicPageMediaResponse    `json:"thumbnail"`
	Gallery   []PublicPageMediaResponse   `json:"gallery"`
	SEO       PublicPageSEOResponse       `json:"seo"`
	Sections  []PublicPageSectionResponse `json:"sections"`
}

type GetPublicPageResponse struct {
	Data PublicPageResponse `json:"data"`
}

type PublicErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
