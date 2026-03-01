package settings_test

import (
	"testing"

	llm_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_agent_reasoning_compute/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentReasoningComputeSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_compute.base_url": "https://llm.example.com/",
		}, "."), nil)

		settings := llm_settings.NewAgentReasoningComputeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://llm.example.com", settings.BaseUrl)
	})
}
