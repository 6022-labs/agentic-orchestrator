package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const MigrationsSettingsKey = "orchestrator_migrations"

type MigrationsSettings struct {
	MigrationsFolderPath string `mapstructure:"folder_path"`
}

func NewMigrationsSettings(
	logger *zap.Logger,
	v *viper.Viper,
) *MigrationsSettings {
	if val, ok := os.LookupEnv("ORCHESTRATOR_MIGRATIONS__FOLDER_PATH"); ok {
		v.Set(MigrationsSettingsKey+".folder_path", val)
	}

	var settings MigrationsSettings
	if err := v.UnmarshalKey(MigrationsSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal migrations settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.MigrationsFolderPath)) == 0 {
		defaultPath := "./migrations"
		logger.Warn("orchestrator_migrations.folder_path is not set, using default path", zap.String("path", defaultPath))
		settings.MigrationsFolderPath = defaultPath
	}

	if _, err := os.Stat(settings.MigrationsFolderPath); os.IsNotExist(err) {
		logger.Fatal("orchestrator_migrations folder path does not exist", zap.String("path", settings.MigrationsFolderPath), zap.Error(err))
	}

	return &settings
}
