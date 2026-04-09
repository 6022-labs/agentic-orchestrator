package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const JwtSettingsKey = "jwt"

type JwtSettings struct {
	Hash       string `koanf:"hash"`
	LivingTime int    `koanf:"living_time"`
}

func NewJwtSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *JwtSettings {
	var settings JwtSettings
	if err := k.Unmarshal(JwtSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal jwt settings", zap.Error(err))
	}

	settings.Hash = strings.TrimSpace(settings.Hash)

	if len(settings.Hash) == 0 {
		logger.Fatal("please set your jwt.hash value in your configuration")
	}

	if settings.LivingTime <= 0 {
		logger.Fatal("please set your jwt.living_time value in your configuration")
	}

	return &settings
}
