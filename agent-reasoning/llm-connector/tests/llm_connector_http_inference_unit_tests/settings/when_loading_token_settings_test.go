package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingTokenSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("token", map[string]any{
			"max_tokens_per_request": uint(1234),
		})

		result := settings.NewTokenSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1234), result.MaxTokensPerRequest)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("TOKEN__MAX_TOKENS_PER_REQUEST", "2048")

		v := viper.New()
		result := settings.NewTokenSettings(zap.NewNop(), v)

		assert.NotNil(t, result)
		assert.Equal(t, uint(2048), result.MaxTokensPerRequest)
	})
}
