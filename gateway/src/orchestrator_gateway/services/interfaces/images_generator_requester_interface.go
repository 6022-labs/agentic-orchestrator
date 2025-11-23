package interfaces

type GenerateImagesRequests struct {
	Name              string
	Role              string
	Creator           string
	EnsDomain         string
	OriginEnsDomain   *string
	OriginGeneration  *uint
	BaseImageBytes    []byte
	BaseImageMimeType string
	Wallets           []string
	Stars             uint8
}

type ImagesGeneratorRequesterInterface interface {
	GenerateImages(request *GenerateImagesRequests) (map[string][]byte, error)
}
