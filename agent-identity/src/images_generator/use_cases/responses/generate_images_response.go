package responses

type GenerateImagesResponse struct {
	Images map[string][]byte `json:"images"`
}
