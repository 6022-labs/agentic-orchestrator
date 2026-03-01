package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute_http_inference/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingXSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"x.base_url": "https://x.example.com",
		}, "."), nil)

		result := settings.NewXSettings(zap.NewNop(), k)

		assert.NotNil(t, result)
		assert.Equal(t, "https://x.example.com", result.BaseUrl)
	})
}
