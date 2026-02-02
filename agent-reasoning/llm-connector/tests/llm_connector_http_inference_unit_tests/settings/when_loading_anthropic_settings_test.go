package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAnthropicSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("anthropic", map[string]any{
			"base_url": "https://anthropic.example.com",
		})

		result := settings.NewAnthropicSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, "https://anthropic.example.com", result.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("ANTHROPIC__BASE_URL", "https://anthropic.env")

		v := viper.New()
		result := settings.NewAnthropicSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, "https://anthropic.env", result.BaseUrl)
	})
}
