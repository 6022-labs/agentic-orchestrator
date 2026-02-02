package settings_test

import (
	"testing"

	platform_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_platform_bridge/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingPlatformBridgesSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("platform_bridges", map[string]any{
			"slack_bridge_url": "https://slack.bridge",
		})

		settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://slack.bridge", settings.SlackBridgeUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("PLATFORM_BRIDGES__SLACK_BRIDGE_URL", "https://slack.env")

		v := viper.New()
		settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://slack.env", settings.SlackBridgeUrl)
	})
}
