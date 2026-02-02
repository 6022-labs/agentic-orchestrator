package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const TokenSettingsKey = "token"

type TokenSettings struct {
	MaxTokensPerRequest uint `mapstructure:"max_tokens_per_request"`
}

func NewTokenSettings(logger *zap.Logger, v *viper.Viper) *TokenSettings {
	v.SetDefault(TokenSettingsKey+".max_tokens_per_request", uint(100000))
	if val, ok := os.LookupEnv("TOKEN__MAX_TOKENS_PER_REQUEST"); ok {
		v.Set(TokenSettingsKey+".max_tokens_per_request", val)
	}

	var settings TokenSettings
	if err := v.UnmarshalKey(TokenSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal token settings", zap.Error(err))
	}

	if settings.MaxTokensPerRequest == 0 {
		logger.Info("token.max_tokens_per_request is not set, using default value", zap.Uint("value", 100000))
		settings.MaxTokensPerRequest = uint(v.GetUint(TokenSettingsKey + ".max_tokens_per_request"))
	}

	return &settings
}
