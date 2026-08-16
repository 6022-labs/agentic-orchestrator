package settings_test

import (
	"testing"

	conversation_broker_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/settings"
	"github.com/google/uuid"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingConversationBrokerSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":     "https://broker.example.com/",
			"conversation_broker.callback_url": "http://slack-bridge.local/broker/callback",
			"conversation_broker.bridge_id":    id.String(),
		}, "."), nil)

		t.Run("Should load them and trim the base url", func(t *testing.T) {
			t.Parallel()

			settings := conversation_broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://broker.example.com", settings.BaseUrl)
			assert.Equal(t, "http://slack-bridge.local/broker/callback", settings.CallbackUrl)
			assert.Equal(t, id, settings.BridgeID)
		})
	})

	t.Run("Given no request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":     "https://broker.example.com",
			"conversation_broker.callback_url": "http://slack-bridge.local/broker/callback",
			"conversation_broker.bridge_id":    uuid.New().String(),
		}, "."), nil)

		t.Run("Should default to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := conversation_broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a zero request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":                "https://broker.example.com",
			"conversation_broker.callback_url":            "http://slack-bridge.local/broker/callback",
			"conversation_broker.bridge_id":               uuid.New().String(),
			"conversation_broker.request_timeout_seconds": 0,
		}, "."), nil)

		t.Run("Should fall back to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := conversation_broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a configured request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":                "https://broker.example.com",
			"conversation_broker.callback_url":            "http://slack-bridge.local/broker/callback",
			"conversation_broker.bridge_id":               uuid.New().String(),
			"conversation_broker.request_timeout_seconds": 5,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := conversation_broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(5), settings.RequestTimeoutSeconds)
		})
	})
}
