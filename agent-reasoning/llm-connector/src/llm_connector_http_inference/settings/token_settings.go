package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const TokenSettingsKey = "token"

type TokenSettings struct {
	MaxTokensPerRequest uint `koanf:"max_tokens_per_request"`
}

func NewTokenSettings(logger *zap.Logger, k *koanf.Koanf) *TokenSettings {
	settings := TokenSettings{MaxTokensPerRequest: uint(100000)}
	if err := k.Unmarshal(TokenSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal token settings", zap.Error(err))
	}

	if settings.MaxTokensPerRequest == 0 {
		logger.Info("token.max_tokens_per_request is not set, using default value", zap.Uint("value", 100000))
		settings.MaxTokensPerRequest = 100000
	}

	return &settings
}
