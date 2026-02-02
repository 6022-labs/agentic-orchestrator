package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const RedisSettingsKey = "redis"

type RedisSettings struct {
	DB       int64  `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

func NewRedisSettings(logger *zap.Logger, v *viper.Viper) *RedisSettings {
	v.SetDefault(RedisSettingsKey+".db", int64(0))

	if val, ok := os.LookupEnv("REDIS__HOST"); ok {
		v.Set(RedisSettingsKey+".host", val)
	}
	if val, ok := os.LookupEnv("REDIS__PASSWORD"); ok {
		v.Set(RedisSettingsKey+".password", val)
	}
	if val, ok := os.LookupEnv("REDIS__DB"); ok {
		v.Set(RedisSettingsKey+".db", val)
	}

	var settings RedisSettings
	if err := v.UnmarshalKey(RedisSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal redis settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.Host)) == 0 {
		logger.Fatal("please set your redis.host value in your configuration")
	}

	if len(strings.TrimSpace(settings.Password)) == 0 {
		logger.Warn("redis password is not set")
	}

	if settings.DB < 0 {
		logger.Warn("invalid redis.db value, defaulting to 0")
		settings.DB = 0
	}

	return &settings
}
