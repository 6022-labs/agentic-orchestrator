package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const JwtSettingsKey = "jwt"

const minHashLength = 32

type JwtSettings struct {
	Hash       string `koanf:"hash"`
	LivingTime int    `koanf:"living_time"`
	NonceTtl   int    `koanf:"nonce_ttl"`
}

func NewJwtSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *JwtSettings {
	settings := JwtSettings{
		NonceTtl: 300,
	}
	if err := k.Unmarshal(JwtSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal jwt settings", zap.Error(err))
	}

	settings.Hash = strings.TrimSpace(settings.Hash)

	if len(settings.Hash) < minHashLength {
		logger.Fatal("jwt.hash must be at least 32 characters")
	}

	if settings.LivingTime <= 0 {
		logger.Fatal("please set your jwt.living_time value in your configuration")
	}

	if settings.NonceTtl <= 0 {
		logger.Warn("invalid jwt.nonce_ttl value, defaulting to 300")
		settings.NonceTtl = 300
	}

	return &settings
}
