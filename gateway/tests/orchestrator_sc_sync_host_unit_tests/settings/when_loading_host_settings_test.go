package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingHostSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"host.use_api":          true,
			"host.use_listeners":    false,
			"host.api_port":         3001,
			"host.app_name":         "SC-Sync",
			"host.body_limit_bytes": 4096,
		}, "."), nil)

		settings := host_settings.NewHostSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.True(t, settings.UseApi)
		assert.False(t, settings.UseListeners)
		assert.Equal(t, 3001, settings.ApiPort)
		assert.Equal(t, "SC-Sync", settings.AppName)
		assert.Equal(t, 4096, settings.BodyLimitBytes)
	})
}
