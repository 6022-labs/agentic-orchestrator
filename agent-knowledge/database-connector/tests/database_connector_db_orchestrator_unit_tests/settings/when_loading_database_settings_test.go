package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorDbOrchestratorDatabaseSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"database.url": "postgres://user:pass@localhost:5432/orchestrator",
		}, "."), nil)

		settings := db_settings.NewDatabaseSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/orchestrator", settings.Url)
	})
}
