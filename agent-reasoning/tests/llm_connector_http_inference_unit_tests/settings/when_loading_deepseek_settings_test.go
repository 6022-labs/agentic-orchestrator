package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingDeepSeekSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"deepseek.base_url": "https://deepseek.example.com",
		}, "."), nil)

		result := settings.NewDeepSeekSettings(zap.NewNop(), k)

		assert.NotNil(t, result)
		assert.Equal(t, "https://deepseek.example.com", result.BaseUrl)
	})
}
