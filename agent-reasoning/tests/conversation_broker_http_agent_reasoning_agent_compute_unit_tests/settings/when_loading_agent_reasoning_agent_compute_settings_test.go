package settings_test

import (
	"testing"

	agent_compute_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_reasoning_agent_compute/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentReasoningAgentComputeSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_reasoning_agent_compute.base_url": "https://agent-compute.example.com",
		}, "."), nil)

		settings := agent_compute_settings.NewAgentReasoningAgentComputeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
	})

	t.Run("Trims trailing slash from base_url", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_reasoning_agent_compute.base_url": "https://agent-compute.example.com/",
		}, "."), nil)

		settings := agent_compute_settings.NewAgentReasoningAgentComputeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
	})
}
