package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeInternalHostSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"host.api_port":         7101,
			"host.app_name":         "slack-bridge-internal",
			"host.body_limit_bytes": 4096,
		}, "."), nil)

		settings := host_settings.NewHostSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 7101, settings.ApiPort)
		assert.Equal(t, "slack-bridge-internal", settings.AppName)
		assert.Equal(t, 4096, settings.BodyLimitBytes)
	})
}
