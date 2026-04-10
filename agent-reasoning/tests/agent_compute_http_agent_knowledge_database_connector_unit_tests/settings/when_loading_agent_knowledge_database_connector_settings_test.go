package settings_test

import (
	"testing"

	agent_knowledge_database_connector_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_agent_knowledge_database_connector/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentKnowledgeDatabaseConnectorSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_knowledge_database_connector.base_url": "https://agent-knowledge-database-connector.example.com",
		}, "."), nil)

		settings := agent_knowledge_database_connector_settings.NewAgentKnowledgeDatabaseConnectorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-knowledge-database-connector.example.com", settings.BaseUrl)
	})

	t.Run("Trims trailing slash from base_url", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_knowledge_database_connector.base_url": "https://agent-knowledge-database-connector.example.com/",
		}, "."), nil)

		settings := agent_knowledge_database_connector_settings.NewAgentKnowledgeDatabaseConnectorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-knowledge-database-connector.example.com", settings.BaseUrl)
	})
}
