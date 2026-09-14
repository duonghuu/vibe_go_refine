package dto

import "time"

type ImageContentIDURI struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type ImageContentFieldRuleRequest struct {
	Enabled  bool `json:"enabled"`
	Required bool `json:"required"`
}

type ImageContentFieldConfigRequest struct {
	Name                 ImageContentFieldRuleRequest `json:"name"`
	Description          ImageContentFieldRuleRequest `json:"description"`
	SecondaryDescription ImageContentFieldRuleRequest `json:"secondaryDescription"`
	URL                  ImageContentFieldRuleRequest `json:"url"`
}

type GetImageContentTypesQuery struct {
	Start  int    `form:"_start" binding:"omitempty,min=0"`
	End    int    `form:"_end" binding:"omitempty,min=1,max=1000"`
	Search string `form:"q" binding:"omitempty,max=100"`
	Status string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	SortBy string `form:"_sort" binding:"omitempty,oneof=id code name status sortOrder createdAt updatedAt"`
	Order  string `form:"_order" binding:"omitempty,oneof=ASC DESC asc desc"`
}

type CreateImageContentTypeRequest struct {
	Code        string                          `json:"code" binding:"required,max=50"`
	Name        string                          `json:"name" binding:"required,max=100"`
	Status      string                          `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
	SortOrder   int                             `json:"sortOrder" binding:"gte=0"`
	MaxItems    *uint                           `json:"maxItems" binding:"omitempty,min=1,max=1000"`
	FieldConfig *ImageContentFieldConfigRequest `json:"fieldConfig" binding:"required"`
}

type UpdateImageContentTypeRequest struct {
	Name        string                          `json:"name" binding:"required,max=100"`
	Status      string                          `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
	SortOrder   int                             `json:"sortOrder" binding:"gte=0"`
	MaxItems    *uint                           `json:"maxItems" binding:"omitempty,min=1,max=1000"`
	FieldConfig *ImageContentFieldConfigRequest `json:"fieldConfig" binding:"required"`
}

type GetImageContentsQuery struct {
	Start    int    `form:"_start" binding:"omitempty,min=0"`
	End      int    `form:"_end" binding:"omitempty,min=1,max=1000"`
	Search   string `form:"q" binding:"omitempty,max=255"`
	TypeCode string `form:"typeCode" binding:"omitempty,max=50"`
	Status   string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	SortBy   string `form:"_sort" binding:"omitempty,oneof=id typeCode name status sortOrder createdAt updatedAt"`
	Order    string `form:"_order" binding:"omitempty,oneof=ASC DESC asc desc"`
}

type CreateImageContentRequest struct {
	TypeCode             string  `json:"typeCode" binding:"required,max=50"`
	MediaID              uint    `json:"mediaId" binding:"required,min=1"`
	Name                 *string `json:"name" binding:"omitempty,max=255"`
	Description          *string `json:"description" binding:"omitempty,max=5000"`
	SecondaryDescription *string `json:"secondaryDescription" binding:"omitempty,max=5000"`
	URL                  *string `json:"url" binding:"omitempty,max=500"`
	SortOrder            *int    `json:"sortOrder" binding:"omitempty,gte=0"`
	Status               string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type UpdateImageContentRequest struct {
	MediaID              uint    `json:"mediaId" binding:"required,min=1"`
	Name                 *string `json:"name" binding:"omitempty,max=255"`
	Description          *string `json:"description" binding:"omitempty,max=5000"`
	SecondaryDescription *string `json:"secondaryDescription" binding:"omitempty,max=5000"`
	URL                  *string `json:"url" binding:"omitempty,max=500"`
	Status               string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type ImageContentOrderItem struct {
	ID        uint `json:"id" binding:"required,min=1"`
	SortOrder int  `json:"sortOrder" binding:"gte=0"`
}

type ReorderImageContentsRequest struct {
	TypeCode string                  `json:"typeCode" binding:"required,max=50"`
	Items    []ImageContentOrderItem `json:"items" binding:"required,max=1000,dive"`
}

type GetPublicImageContentsQuery struct {
	TypeCodes string `form:"typeCodes" binding:"omitempty,max=1024"`
}

type ImageContentFieldRuleResponse struct {
	Enabled  bool `json:"enabled"`
	Required bool `json:"required"`
}

type ImageContentFieldConfigResponse struct {
	Name                 ImageContentFieldRuleResponse `json:"name"`
	Description          ImageContentFieldRuleResponse `json:"description"`
	SecondaryDescription ImageContentFieldRuleResponse `json:"secondaryDescription"`
	URL                  ImageContentFieldRuleResponse `json:"url"`
}

type ImageContentTypeResponse struct {
	ID          uint                            `json:"id"`
	Code        string                          `json:"code"`
	Name        string                          `json:"name"`
	Status      string                          `json:"status"`
	SortOrder   int                             `json:"sortOrder"`
	MaxItems    *uint                           `json:"maxItems"`
	FieldConfig ImageContentFieldConfigResponse `json:"fieldConfig"`
	ItemCount   int64                           `json:"itemCount"`
	CreatedAt   time.Time                       `json:"createdAt"`
	UpdatedAt   time.Time                       `json:"updatedAt"`
}

type AdminImageContentMediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
	MimeType     string `json:"mimeType"`
	Status       string `json:"status"`
}

type ImageContentResponse struct {
	ID                   uint                            `json:"id"`
	TypeCode             string                          `json:"typeCode"`
	TypeName             string                          `json:"typeName"`
	FieldConfig          ImageContentFieldConfigResponse `json:"fieldConfig"`
	MediaID              uint                            `json:"mediaId"`
	Media                AdminImageContentMediaResponse  `json:"media"`
	Name                 *string                         `json:"name"`
	Description          *string                         `json:"description"`
	SecondaryDescription *string                         `json:"secondaryDescription"`
	URL                  *string                         `json:"url"`
	SortOrder            int                             `json:"sortOrder"`
	Status               string                          `json:"status"`
	CreatedBy            uint                            `json:"createdBy"`
	CreatedAt            time.Time                       `json:"createdAt"`
	UpdatedAt            time.Time                       `json:"updatedAt"`
}

type ImageContentTypeListResponse struct {
	Data  []ImageContentTypeResponse `json:"data"`
	Total int64                      `json:"total"`
}

type ImageContentListResponse struct {
	Data  []ImageContentResponse `json:"data"`
	Total int64                  `json:"total"`
}

type ImageContentErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type PublicImageContentMediaResponse struct {
	ID           uint   `json:"id"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	MediumURL    string `json:"mediumUrl,omitempty"`
	MimeType     string `json:"mimeType"`
}

type PublicImageContentItemResponse struct {
	ID                   uint                            `json:"id"`
	Name                 *string                         `json:"name,omitempty"`
	Description          *string                         `json:"description,omitempty"`
	SecondaryDescription *string                         `json:"secondaryDescription,omitempty"`
	URL                  *string                         `json:"url,omitempty"`
	SortOrder            int                             `json:"sortOrder"`
	Media                PublicImageContentMediaResponse `json:"media"`
}

type PublicImageContentGroupResponse struct {
	TypeCode string                           `json:"typeCode"`
	Name     string                           `json:"name"`
	Items    []PublicImageContentItemResponse `json:"items"`
}

type PublicImageContentResponse struct {
	Data []PublicImageContentGroupResponse `json:"data"`
}
