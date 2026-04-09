package settings

import (
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const ImageGenerationSettingsKey = "image_generation"

type ImageGenerationSettings struct {
	ImageMarkPath *string `koanf:"image_mark_path"`
}

func NewImageGenerationSettings(logger *zap.Logger, k *koanf.Koanf) *ImageGenerationSettings {
	var settings ImageGenerationSettings
	if err := k.Unmarshal(ImageGenerationSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal image generation settings", zap.Error(err))
	}

	if settings.ImageMarkPath != nil {
		value := strings.TrimSpace(*settings.ImageMarkPath)
		if len(value) == 0 {
			settings.ImageMarkPath = nil
		} else if _, err := os.Stat(value); os.IsNotExist(err) {
			logger.Warn("The file specified in image_generation.image_mark_path does not exist", zap.String("path", value))
			settings.ImageMarkPath = nil
		} else {
			settings.ImageMarkPath = &value
		}
	}

	if settings.ImageMarkPath == nil {
		logger.Info("image_generation.image_mark_path is not set, nft generation will not be marked")
	}

	return &settings
}
