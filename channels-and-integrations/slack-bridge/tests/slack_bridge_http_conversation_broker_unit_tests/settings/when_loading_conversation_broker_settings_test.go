package settings_test

import (
	"testing"

	conversation_broker_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/settings"
	"github.com/google/uuid"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingConversationBrokerSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		id := uuid.New()
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":     "https://broker.example.com/",
			"conversation_broker.callback_url": "http://slack-bridge.local/internal/broker/callback",
			"conversation_broker.bridge_id":    id.String(),
		}, "."), nil)

		settings := conversation_broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://broker.example.com", settings.BaseUrl)
		assert.Equal(t, "http://slack-bridge.local/internal/broker/callback", settings.CallbackUrl)
		assert.Equal(t, id, settings.BridgeID)
	})
}
