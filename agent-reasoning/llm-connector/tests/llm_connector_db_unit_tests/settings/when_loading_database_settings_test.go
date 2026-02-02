package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingLlmConnectorDatabaseSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("database", map[string]any{
			"url": "postgres://user:pass@localhost:5432/llm",
		})

		settings := db_settings.NewOrchestratorDatabaseSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/llm", settings.Url)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("DATABASE__URL", "postgres://user:pass@localhost:5432/llm-env")

		v := viper.New()
		settings := db_settings.NewOrchestratorDatabaseSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/llm-env", settings.Url)
	})
}
