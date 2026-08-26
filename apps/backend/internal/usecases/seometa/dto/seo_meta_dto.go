package dto

import "encoding/json"

type EntityURI struct {
	EntityType string `uri:"entityType" binding:"required,max=50"`
	EntityID   uint   `uri:"entityId" binding:"required,min=1"`
}

type UpsertRequest struct {
	MetaTitle          *string         `json:"meta_title"`
	MetaDescription    *string         `json:"meta_description"`
	MetaKeywords       *string         `json:"meta_keywords"`
	CanonicalURL       *string         `json:"canonical_url"`
	OGTitle            *string         `json:"og_title"`
	OGDescription      *string         `json:"og_description"`
	OGImage            *string         `json:"og_image"`
	OGType             *string         `json:"og_type"`
	TwitterTitle       *string         `json:"twitter_title"`
	TwitterDescription *string         `json:"twitter_description"`
	TwitterImage       *string         `json:"twitter_image"`
	TwitterCard        *string         `json:"twitter_card"`
	Robots             *string         `json:"robots"`
	SchemaJSON         json.RawMessage `json:"schema_json"`
}

type SEOMetaResponse struct {
	EntityType           string          `json:"entity_type"`
	EntityID             uint            `json:"entity_id"`
	MetaTitle            *string         `json:"meta_title"`
	ResolvedTitle        string          `json:"resolved_title"`
	MetaDescription      *string         `json:"meta_description"`
	ResolvedDescription  string          `json:"resolved_description"`
	MetaKeywords         *string         `json:"meta_keywords"`
	CanonicalURL         *string         `json:"canonical_url"`
	ResolvedCanonicalURL string          `json:"resolved_canonical_url"`
	OGTitle              *string         `json:"og_title"`
	OGDescription        *string         `json:"og_description"`
	OGImage              *string         `json:"og_image"`
	OGType               *string         `json:"og_type"`
	TwitterTitle         *string         `json:"twitter_title"`
	TwitterDescription   *string         `json:"twitter_description"`
	TwitterImage         *string         `json:"twitter_image"`
	TwitterCard          *string         `json:"twitter_card"`
	Robots               string          `json:"robots"`
	SchemaJSON           json.RawMessage `json:"schema_json"`
}

type Response struct {
	Data SEOMetaResponse `json:"data"`
}
type UpsertResponse struct {
	Data    SEOMetaResponse `json:"data"`
	Message string          `json:"message"`
}
type DeleteResponse struct {
	Message string `json:"message"`
}
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
