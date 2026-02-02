package settings_test

import (
	"testing"

	llm_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_reasoning_llm_connector/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentLlmConnectorSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("agent_llm_connector", map[string]any{
			"base_url": "https://llm.example.com/",
		})

		settings := llm_settings.NewAgentLLMConnectorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://llm.example.com", settings.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("AGENT_LLM_CONNECTOR__BASE_URL", "https://llm.env/")

		v := viper.New()
		settings := llm_settings.NewAgentLLMConnectorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://llm.env", settings.BaseUrl)
	})
}
