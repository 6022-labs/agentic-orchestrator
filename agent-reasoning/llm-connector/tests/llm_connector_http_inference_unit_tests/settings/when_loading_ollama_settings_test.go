package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_http_inference/settings"
	"github.com/stretchr/testify/assert"
)

func TestWhenLoadingOllamaSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		result := settings.NewOllamaSettings()

		assert.NotNil(t, result)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		result := settings.NewOllamaSettings()

		assert.NotNil(t, result)
	})
}
