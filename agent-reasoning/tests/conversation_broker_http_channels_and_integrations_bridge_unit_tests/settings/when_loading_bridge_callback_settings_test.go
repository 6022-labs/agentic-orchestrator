package settings_test

import (
	"testing"

	bridge_callback_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_channels_and_integrations_bridge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingBridgeCallbackSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"bridge_callback.request_timeout_seconds": 60,
		}, "."), nil)

		settings := bridge_callback_settings.NewBridgeCallbackSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 60, settings.RequestTimeoutSeconds)
	})

	t.Run("Defaults request_timeout_seconds when unset", func(t *testing.T) {
		k := koanf.New(".")

		settings := bridge_callback_settings.NewBridgeCallbackSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 30, settings.RequestTimeoutSeconds)
	})
}
