package settings_test

import (
	"testing"

	openai_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingOpenAISettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("openai", map[string]any{
			"base_url": "https://openai.example.com",
		})

		settings := openai_settings.NewOpenAISettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://openai.example.com", settings.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("OPENAI__BASE_URL", "https://openai.env")

		v := viper.New()
		settings := openai_settings.NewOpenAISettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://openai.env", settings.BaseUrl)
	})
}
