package entity

// ImageContentStatus represents the lifecycle state shared by image content
// types and items.
type ImageContentStatus string

const (
	ImageContentStatusActive   ImageContentStatus = "ACTIVE"
	ImageContentStatusInactive ImageContentStatus = "INACTIVE"
)

// ImageContentFieldRule controls whether a metadata field is available and
// whether the field must be populated for an image content item.
type ImageContentFieldRule struct {
	Enabled  bool `json:"enabled"`
	Required bool `json:"required"`
}

// ImageContentFieldConfig is intentionally closed over the supported metadata
// fields. Adding a new field requires an explicit backend/schema change.
type ImageContentFieldConfig struct {
	Name                 ImageContentFieldRule `json:"name"`
	Description          ImageContentFieldRule `json:"description"`
	SecondaryDescription ImageContentFieldRule `json:"secondaryDescription"`
	URL                  ImageContentFieldRule `json:"url"`
}
