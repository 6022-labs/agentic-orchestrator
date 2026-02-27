package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_db/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingLlmConnectorDatabaseSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"database.url": "postgres://user:pass@localhost:5432/llm",
		}, "."), nil)

		settings := db_settings.NewOrchestratorDatabaseSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "postgres://user:pass@localhost:5432/llm", settings.Url)
	})
}
