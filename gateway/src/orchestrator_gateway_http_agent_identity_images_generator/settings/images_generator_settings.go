package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type ImagesGeneratorSettings struct {
	BaseUrl string
}

func NewImagesGeneratorSettings(logger *zap.Logger) *ImagesGeneratorSettings {
	imageGeneratorBaseUrl := os.Getenv("IMAGE_GENERATOR_BASE_URL")

	if len(strings.TrimSpace(imageGeneratorBaseUrl)) == 0 {
		logger.Fatal("IMAGE_GENERATOR_BASE_URL is not set")
	}

	imageGeneratorBaseUrl = strings.TrimRight(imageGeneratorBaseUrl, "/")

	return &ImagesGeneratorSettings{
		BaseUrl: imageGeneratorBaseUrl,
	}
}
