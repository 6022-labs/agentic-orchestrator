package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ImagesGeneratorSettingsKey = "images_generator"

type ImagesGeneratorSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewImagesGeneratorSettings(logger *zap.Logger, k *koanf.Koanf) *ImagesGeneratorSettings {
	var settings ImagesGeneratorSettings
	if err := k.Unmarshal(ImagesGeneratorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal images generator settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)
	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your images_generator.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")
	return &settings
}
