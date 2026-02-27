package settings_test

import (
	"testing"

	platform_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_platform_bridge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingPlatformBridgesSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"platform_bridges.slack_bridge_url": "https://slack.bridge",
		}, "."), nil)

		settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://slack.bridge", settings.SlackBridgeUrl)
	})
}
