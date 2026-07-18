package settings_test

import (
	"testing"

	broker_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingConversationBrokerSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":                "https://broker.example.com",
			"conversation_broker.request_timeout_seconds": uint64(15),
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://broker.example.com", settings.BaseUrl)
			assert.Equal(t, uint64(15), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a base url with a trailing slash and whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url": "  https://broker.example.com/  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.Equal(t, "https://broker.example.com", settings.BaseUrl)
		})
	})

	t.Run("Given no request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url": "https://broker.example.com",
		}, "."), nil)

		t.Run("Should fall back to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a zero request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"conversation_broker.base_url":                "https://broker.example.com",
			"conversation_broker.request_timeout_seconds": uint64(0),
		}, "."), nil)

		t.Run("Should fall back to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewConversationBrokerSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})
}
