package settings_test

import (
	"testing"

	orchestrator_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeHttpOrchestratorSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"orchestrator.api_base_url": "https://orchestrator.example.com/",
		}, "."), nil)

		settings := orchestrator_settings.NewOrchestratorSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://orchestrator.example.com", settings.ApiBaseUrl)
	})
}
