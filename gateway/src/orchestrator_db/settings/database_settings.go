package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const DatabaseSettingsKey = "database"

type DatabaseSettings struct {
	Url string `mapstructure:"url"`
}

func NewDatabaseSettings(
	logger *zap.Logger,
	v *viper.Viper,
) *DatabaseSettings {
	if val, ok := os.LookupEnv("DATABASE__URL"); ok {
		v.Set(DatabaseSettingsKey+".url", val)
	}

	var settings DatabaseSettings
	if err := v.UnmarshalKey(DatabaseSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal database settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.Url)) == 0 {
		logger.Fatal("please set your database.url value in your configuration")
	}

	return &settings
}
