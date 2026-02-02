package settings_test

import (
	"testing"

	orchestrator_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeHttpOrchestratorSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("orchestrator", map[string]any{
			"api_base_url": "https://orchestrator.example.com/",
		})

		settings := orchestrator_settings.NewOrchestratorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://orchestrator.example.com", settings.ApiBaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("ORCHESTRATOR__API_BASE_URL", "https://orchestrator.env/")

		v := viper.New()
		settings := orchestrator_settings.NewOrchestratorSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://orchestrator.env", settings.ApiBaseUrl)
	})
}
