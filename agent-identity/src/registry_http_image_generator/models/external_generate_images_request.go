package models

type ExternalGenerateImagesRequest struct {
	Name              string   `form:"name"`
	Role              string   `form:"role"`
	Creator           string   `form:"creator"`
	EnsDomain         string   `form:"ensDomain"`
	OriginEnsDomain   *string  `form:"originEnsDomain"`
	OriginGeneration  *uint    `form:"originGeneration"`
	BaseImageBytes    []byte   `form:"baseImageBytes"`
	BaseImageMimeType string   `form:"baseImageMimeType"`
	Wallets           []string `form:"wallets"`
	Stars             uint8    `form:"stars"`
}
