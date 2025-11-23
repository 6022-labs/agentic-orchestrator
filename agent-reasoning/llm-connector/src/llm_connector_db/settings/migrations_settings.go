package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type MigrationsSettings struct {
	MigrationsFolderPath string
}

func NewMigrationsSettings(
	logger *zap.Logger,
) *MigrationsSettings {
	migrationsFolderPath := os.Getenv("MIGRATIONS_FOLDER_PATH")
	if len(strings.TrimSpace(migrationsFolderPath)) == 0 {
		logger.Fatal("MIGRATIONS_FOLDER_PATH environment variable is required but not set")
	}

	return &MigrationsSettings{
		MigrationsFolderPath: migrationsFolderPath,
	}
}
