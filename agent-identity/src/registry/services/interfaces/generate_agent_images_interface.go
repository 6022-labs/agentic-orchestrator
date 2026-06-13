package interfaces

import "time"

// GenerateAgentImagesInterface renders an agent's framed images, keyed "default" and "icon" as the collection contract expects.
type GenerateAgentImagesInterface interface {
	Generate(params GenerateImageParams) (map[string][]byte, error)
}

type GenerateImageParams struct {
	Name              string
	EnsDomain         string
	Role              string
	Creator           string
	OriginEnsDomain   *string
	OriginGeneration  *uint
	BaseImageBytes    []byte
	BaseImageMimeType string
	Wallets           []string
	CreatedAt         time.Time
	Stars             uint8 // capped at 8
}
