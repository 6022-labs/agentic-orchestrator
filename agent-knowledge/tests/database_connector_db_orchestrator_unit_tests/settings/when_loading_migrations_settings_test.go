package settings_test

import (
	"testing"

	migrations_settings "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_orchestrator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorDbOrchestratorMigrationsSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		folder := t.TempDir()
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"orchestrator_migrations.folder_path": folder,
		}, "."), nil)

		settings := migrations_settings.NewMigrationsSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, folder, settings.MigrationsFolderPath)
	})
}
