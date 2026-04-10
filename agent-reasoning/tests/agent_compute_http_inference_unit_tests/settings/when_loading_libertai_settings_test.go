package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingLibertAISettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"libertai.base_url": "https://libertai.example.com",
		}, "."), nil)

		result := settings.NewLibertAISettings(zap.NewNop(), k)

		assert.NotNil(t, result)
		assert.Equal(t, "https://libertai.example.com", result.BaseUrl)
	})
}
