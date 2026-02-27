package settings

import (
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const MigrationsSettingsKey = "agent_migrations"

type MigrationsSettings struct {
	MigrationsFolderPath string `koanf:"folder_path"`
}

func NewMigrationsSettings(
	logger *zap.Logger,
	k *koanf.Koanf,
) *MigrationsSettings {
	var settings MigrationsSettings
	if err := k.Unmarshal(MigrationsSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal migrations settings", zap.Error(err))
	}

	settings.MigrationsFolderPath = strings.TrimSpace(settings.MigrationsFolderPath)

	if len(settings.MigrationsFolderPath) == 0 {
		defaultPath := "./migrations"
		logger.Warn("agent_migrations.folder_path is not set, using default path", zap.String("path", defaultPath))
		settings.MigrationsFolderPath = defaultPath
	}

	if _, err := os.Stat(settings.MigrationsFolderPath); os.IsNotExist(err) {
		logger.Fatal("agent_migrations folder path does not exist", zap.String("path", settings.MigrationsFolderPath), zap.Error(err))
	}

	return &settings
}
