package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const JwtSettingsKey = "jwt"

const minHashLength = 32

const (
	defaultNonceTtl       = 300
	defaultRefreshTtl     = 604800
	defaultRefreshIdleTtl = 86400
)

type JwtSettings struct {
	Hash           string `koanf:"hash"`
	LivingTime     int    `koanf:"living_time"`
	NonceTtl       int    `koanf:"nonce_ttl"`
	RefreshTtl     int    `koanf:"refresh_ttl"`
	RefreshIdleTtl int    `koanf:"refresh_idle_ttl"`
}

func NewJwtSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *JwtSettings {
	settings := JwtSettings{
		NonceTtl:       defaultNonceTtl,
		RefreshTtl:     defaultRefreshTtl,
		RefreshIdleTtl: defaultRefreshIdleTtl,
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
		settings.NonceTtl = defaultNonceTtl
	}

	if settings.RefreshTtl <= 0 {
		logger.Warn("invalid jwt.refresh_ttl value, defaulting to 604800")
		settings.RefreshTtl = defaultRefreshTtl
	}

	if settings.RefreshIdleTtl <= 0 {
		logger.Warn("invalid jwt.refresh_idle_ttl value, defaulting to 86400")
		settings.RefreshIdleTtl = defaultRefreshIdleTtl
	}

	if settings.RefreshIdleTtl < 2*settings.LivingTime {
		logger.Warn("jwt.refresh_idle_ttl is under twice jwt.living_time, raising it")
		settings.RefreshIdleTtl = 2 * settings.LivingTime
	}

	return &settings
}
