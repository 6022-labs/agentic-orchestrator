package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ImageGenerationSettingsKey = "image_generation"

type ImageGenerationSettings struct {
	ImageMarkPath *string `mapstructure:"image_mark_path"`
}

func NewImageGenerationSettings(logger *zap.Logger, v *viper.Viper) *ImageGenerationSettings {
	if val, ok := os.LookupEnv("IMAGE_GENERATION__IMAGE_MARK_PATH"); ok {
		v.Set(ImageGenerationSettingsKey+".image_mark_path", val)
	}

	var imageMarkPath *string = nil
	imageMarkPathEnv := v.GetString(ImageGenerationSettingsKey + ".image_mark_path")

	if len(strings.TrimSpace(imageMarkPathEnv)) != 0 {
		if _, err := os.Stat(imageMarkPathEnv); os.IsNotExist(err) {
			logger.Warn("The file specified in IMAGE_GENERATION__IMAGE_MARK_PATH does not exist", zap.String("path", imageMarkPathEnv))
		} else {
			imageMarkPath = &imageMarkPathEnv
		}
	}

	if imageMarkPath == nil {
		logger.Info("image_generation.image_mark_path is not set, nft generation will not be marked")
	}

	return &ImageGenerationSettings{
		ImageMarkPath: imageMarkPath,
	}
}
