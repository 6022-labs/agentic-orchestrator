package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingMistralSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("mistral", map[string]any{
			"base_url": "https://mistral.example.com",
		})

		result := settings.NewMistralSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, "https://mistral.example.com", result.BaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("MISTRAL__BASE_URL", "https://mistral.env")

		v := viper.New()
		result := settings.NewMistralSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, "https://mistral.env", result.BaseUrl)
	})
}
