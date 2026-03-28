package settings_test

import (
	"testing"

	agent_knowledge_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_agent_knowledge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentKnowledgeSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_knowledge.base_url": "https://agent-knowledge.example.com",
		}, "."), nil)

		settings := agent_knowledge_settings.NewAgentKnowledgeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-knowledge.example.com", settings.BaseUrl)
	})

	t.Run("Trims trailing slash from base_url", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_knowledge.base_url": "https://agent-knowledge.example.com/",
		}, "."), nil)

		settings := agent_knowledge_settings.NewAgentKnowledgeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-knowledge.example.com", settings.BaseUrl)
	})
}
