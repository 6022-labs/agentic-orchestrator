package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_host/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingHostSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("host", map[string]any{
			"api_port":         5001,
			"app_name":         "Gateway",
			"body_limit_bytes": 1024,
		})

		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 5001, settings.ApiPort)
		assert.Equal(t, "Gateway", settings.AppName)
		assert.Equal(t, 1024, settings.BodyLimitBytes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("HOST__API_PORT", "5002")
		t.Setenv("HOST__APP_NAME", "GatewayEnv")
		t.Setenv("HOST__BODY_LIMIT_BYTES", "2048")

		v := viper.New()
		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, 5002, settings.ApiPort)
		assert.Equal(t, "GatewayEnv", settings.AppName)
		assert.Equal(t, 2048, settings.BodyLimitBytes)
	})
}
