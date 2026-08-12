package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRpcNodeSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.http_url":                "https://example.com/rpc",
			"rpc_node.request_timeout_seconds": 30,
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://example.com/rpc", settings.HttpUrl)
			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given no configured timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.http_url": "https://example.com/rpc",
		}, "."), nil)

		t.Run("Should fall back to the default", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(10), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a url with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.http_url": "  https://example.com/rpc  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.Equal(t, "https://example.com/rpc", settings.HttpUrl)
		})
	})
}
