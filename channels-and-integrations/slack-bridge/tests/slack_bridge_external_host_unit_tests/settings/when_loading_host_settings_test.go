package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeExternalHostSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("host", map[string]any{
			"api_port":         7001,
			"app_name":         "slack-bridge-external",
			"body_limit_bytes": 1024,
		})

		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 7001, settings.ApiPort)
		assert.Equal(t, "slack-bridge-external", settings.AppName)
		assert.Equal(t, 1024, settings.BodyLimitBytes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("HOST__API_PORT", "7002")
		t.Setenv("HOST__APP_NAME", "slack-bridge-external")
		t.Setenv("HOST__BODY_LIMIT_BYTES", "2048")

		v := viper.New()
		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 7002, settings.ApiPort)
		assert.Equal(t, "slack-bridge-external", settings.AppName)
		assert.Equal(t, 2048, settings.BodyLimitBytes)
	})
}
