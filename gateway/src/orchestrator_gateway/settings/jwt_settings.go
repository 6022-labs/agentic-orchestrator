package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const JwtSettingsKey = "jwt"

type JwtSettings struct {
	Hash       string `mapstructure:"hash"`
	LivingTime int    `mapstructure:"living_time"`
}

func NewJwtSettings(
	logger *zap.Logger,
	v *viper.Viper,
) *JwtSettings {
	if val, ok := os.LookupEnv("JWT__HASH"); ok {
		v.Set(JwtSettingsKey+".hash", val)
	}
	if val, ok := os.LookupEnv("JWT__LIVING_TIME"); ok {
		v.Set(JwtSettingsKey+".living_time", val)
	}

	var settings JwtSettings
	if err := v.UnmarshalKey(JwtSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal jwt settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.Hash)) == 0 {
		logger.Fatal("please set your jwt.hash value in your configuration")
	}

	if settings.LivingTime <= 0 {
		logger.Fatal("please set your jwt.living_time value in your configuration")
	}

	return &settings
}
