package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync_host/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingHostSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("host", map[string]any{
			"use_api":          true,
			"use_listeners":    false,
			"api_port":         3001,
			"app_name":         "SC-Sync",
			"body_limit_bytes": 4096,
		})

		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.True(t, settings.UseApi)
		assert.False(t, settings.UseListeners)
		assert.Equal(t, 3001, settings.ApiPort)
		assert.Equal(t, "SC-Sync", settings.AppName)
		assert.Equal(t, 4096, settings.BodyLimitBytes)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("HOST__USE_API", "false")
		t.Setenv("HOST__USE_LISTENERS", "true")
		t.Setenv("HOST__API_PORT", "3002")
		t.Setenv("HOST__APP_NAME", "SC-Sync-Env")
		t.Setenv("HOST__BODY_LIMIT_BYTES", "8192")

		v := viper.New()
		settings := host_settings.NewHostSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.False(t, settings.UseApi)
		assert.True(t, settings.UseListeners)
		assert.Equal(t, 3002, settings.ApiPort)
		assert.Equal(t, "SC-Sync-Env", settings.AppName)
		assert.Equal(t, 8192, settings.BodyLimitBytes)
	})
}
