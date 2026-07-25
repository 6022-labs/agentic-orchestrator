package settings_test

import (
	"testing"

	bridge_callback_settings "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingBridgeCallbackSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"bridge_callback.request_timeout_seconds": 60,
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := bridge_callback_settings.NewBridgeCallbackSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, 60, settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to defaults", func(t *testing.T) {
			t.Parallel()

			settings := bridge_callback_settings.NewBridgeCallbackSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, 30, settings.RequestTimeoutSeconds)
		})
	})
}
