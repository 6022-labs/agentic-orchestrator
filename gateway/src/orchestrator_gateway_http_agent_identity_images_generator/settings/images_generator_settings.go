package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ImagesGeneratorSettingsKey = "images_generator"

type ImagesGeneratorSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewImagesGeneratorSettings(logger *zap.Logger, v *viper.Viper) *ImagesGeneratorSettings {
	if val, ok := os.LookupEnv("IMAGES_GENERATOR__BASE_URL"); ok {
		v.Set(ImagesGeneratorSettingsKey+".base_url", val)
	}

	var settings ImagesGeneratorSettings
	if err := v.UnmarshalKey(ImagesGeneratorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal images generator settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.BaseUrl)) == 0 {
		logger.Fatal("please set your images_generator.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
