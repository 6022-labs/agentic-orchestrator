package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAnthropicSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"anthropic.base_url": "https://anthropic.example.com",
		}, "."), nil)

		result := settings.NewAnthropicSettings(zap.NewNop(), k)

		assert.NotNil(t, result)
		assert.Equal(t, "https://anthropic.example.com", result.BaseUrl)
	})
}
