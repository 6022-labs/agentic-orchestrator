package settings_test

import (
	"testing"

	migrations_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingLlmConnectorMigrationsSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		folder := t.TempDir()

		v := viper.New()
		v.Set("migrations", map[string]any{
			"folder_path": folder,
		})

		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, folder, settings.MigrationsFolderPath)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		folder := t.TempDir()
		t.Setenv("MIGRATIONS__FOLDER_PATH", folder)

		v := viper.New()
		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, folder, settings.MigrationsFolderPath)
	})
}
