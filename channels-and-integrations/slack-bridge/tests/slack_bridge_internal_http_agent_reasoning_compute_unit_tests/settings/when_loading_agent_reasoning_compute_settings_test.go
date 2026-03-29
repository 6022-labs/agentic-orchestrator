package settings_test

import (
	"testing"

	agent_reasoning_compute_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_http_agent_reasoning_compute/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentReasoningComputeSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_compute.base_url": "https://agent-compute.example.com/",
		}, "."), nil)

		settings := agent_reasoning_compute_settings.NewAgentReasoningComputeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
	})
}
