package settings_test

import (
	"testing"

	migrations_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorDbOrchestratorMigrationsSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		folder := t.TempDir()

		v := viper.New()
		v.Set("orchestrator_migrations", map[string]any{
			"folder_path": folder,
		})

		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, folder, settings.MigrationsFolderPath)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		folder := t.TempDir()
		t.Setenv("ORCHESTRATOR_MIGRATIONS__FOLDER_PATH", folder)

		v := viper.New()
		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, folder, settings.MigrationsFolderPath)
	})
}
