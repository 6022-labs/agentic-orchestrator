package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const DatabaseSettingsKey = "database"

type DatabaseSettings struct {
	URL          string `koanf:"url"`
	MaxOpenConns int    `koanf:"max_open_conns"`
}

func NewDatabaseSettings(logger *zap.Logger, k *koanf.Koanf) *DatabaseSettings {
	settings := DatabaseSettings{MaxOpenConns: 20}

	if err := k.Unmarshal(DatabaseSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal database settings", zap.Error(err))
	}

	settings.URL = strings.TrimSpace(settings.URL)
	if settings.URL == "" {
		logger.Fatal("please set your database.url value in your configuration")
	}

	if settings.MaxOpenConns <= 0 {
		logger.Warn("invalid database.max_open_conns value, defaulting to 20")
		settings.MaxOpenConns = 20
	}

	return &settings
}
