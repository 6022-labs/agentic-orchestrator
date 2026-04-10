package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute_http_inference/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingTokenSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"token.max_tokens_per_request": uint(1234),
		}, "."), nil)

		result := settings.NewTokenSettings(zap.NewNop(), k)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1234), result.MaxTokensPerRequest)
	})
}
