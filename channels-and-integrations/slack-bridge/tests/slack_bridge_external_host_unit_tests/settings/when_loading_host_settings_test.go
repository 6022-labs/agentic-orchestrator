package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeExternalHostSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"host.api_port":         7001,
			"host.app_name":         "slack-bridge-external",
			"host.body_limit_bytes": 1024,
		}, "."), nil)

		settings := host_settings.NewHostSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 7001, settings.ApiPort)
		assert.Equal(t, "slack-bridge-external", settings.AppName)
		assert.Equal(t, 1024, settings.BodyLimitBytes)
	})
}
