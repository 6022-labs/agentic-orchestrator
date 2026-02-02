package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDatabaseConnectorDbOrchestratorDatabaseSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("database", map[string]any{
			"url": "postgres://user:pass@localhost:5432/orchestrator",
		})

		settings := db_settings.NewDatabaseSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/orchestrator", settings.Url)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("DATABASE__URL", "postgres://user:pass@localhost:5432/orchestrator-env")

		v := viper.New()
		settings := db_settings.NewDatabaseSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/orchestrator-env", settings.Url)
	})
}
