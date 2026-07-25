package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const DatabaseSettingsKey = "database"

type DatabaseSettings struct {
	Url string `koanf:"url"`
}

func NewConversationBrokerDatabaseSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *DatabaseSettings {
	var settings DatabaseSettings
	if err := k.Unmarshal(DatabaseSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal database settings", zap.Error(err))
	}

	settings.Url = strings.TrimSpace(settings.Url)

	if len(settings.Url) == 0 {
		logger.Fatal("please set your database.url value in your configuration")
	}

	return &settings
}
