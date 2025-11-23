package models

type ExternalGenerateImagesResponse struct {
	Images map[string][]byte `json:"images"`
}
