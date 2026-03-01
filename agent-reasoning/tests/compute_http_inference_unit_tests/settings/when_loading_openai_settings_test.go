package settings_test

import (
	"testing"

	openai_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_inference/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingOpenAISettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"openai.base_url": "https://openai.example.com",
		}, "."), nil)

		settings := openai_settings.NewOpenAISettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://openai.example.com", settings.BaseUrl)
	})
}
