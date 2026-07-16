package dto

type UploadMediaResponse struct {
	ID           uint   `json:"id"`
	FileName     string `json:"fileName"`
	OriginalUrl  string `json:"originalUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	MediumUrl    string `json:"mediumUrl"`
	Status       string `json:"status"`
}
