package settings_test

import (
	"testing"

	settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector_http_inference/settings"
	"github.com/stretchr/testify/assert"
)

func TestWhenLoadingAmazonBedrockSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		result := settings.NewAmazonBedrockSettings()

		assert.NotNil(t, result)
	})
}
