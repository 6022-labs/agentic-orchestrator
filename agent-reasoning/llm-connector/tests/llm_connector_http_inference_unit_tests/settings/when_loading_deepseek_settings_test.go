package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDeepSeekSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("deepseek", map[string]any{
			"base_url": "https://deepseek.example.com",
		})

		result := settings.NewDeepSeekSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, "https://deepseek.example.com", result.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("DEEPSEEK__BASE_URL", "https://deepseek.env")

		v := viper.New()
		result := settings.NewDeepSeekSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, "https://deepseek.env", result.BaseUrl)
	})
}
