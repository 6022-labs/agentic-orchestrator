package settings_test

import (
	"testing"

	migrations_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_db/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingMigrationsSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		tempDir := t.TempDir()
		v := viper.New()
		v.Set("migrations", map[string]any{
			"folder_path": tempDir,
		})

		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, tempDir, settings.MigrationsFolderPath)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		tempDir := t.TempDir()
		t.Setenv("MIGRATIONS__FOLDER_PATH", tempDir)

		v := viper.New()
		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, tempDir, settings.MigrationsFolderPath)
	})
}
