package settings_test

import (
	"testing"

	db_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_knowledge_database_connector/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentKnowledgeDatabaseConnectorSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_knowledge_database_connector.base_url": "https://db.example.com/",
		}, "."), nil)

		settings := db_settings.NewAgentKnowledgeDatabaseConnectorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://db.example.com", settings.BaseUrl)
	})
}
