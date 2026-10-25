package settings_test

import (
	"testing"

	rpc_node_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
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
			"rpc_node.ws_url":                  "wss://example.com/ws",
			"rpc_node.http_url":                "https://example.com/rpc",
			"rpc_node.request_timeout_seconds": 30,
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := rpc_node_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "wss://example.com/ws", settings.WsUrl)
			assert.Equal(t, "https://example.com/rpc", settings.HttpUrl)
			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given only the urls are configured", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.ws_url":   "wss://example.com/ws",
			"rpc_node.http_url": "https://example.com/rpc",
		}, "."), nil)

		t.Run("Should fall back to the default", func(t *testing.T) {
			t.Parallel()

			settings := rpc_node_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(10), settings.RequestTimeoutSeconds)
		})

		t.Run("Should fall back to 10 addresses per log query", func(t *testing.T) {
			t.Parallel()

			settings := rpc_node_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.Equal(t, 10, settings.MaxAddressesPerLogQuery)
		})
	})

	t.Run("Given a configured address limit per log query", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.ws_url":                      "wss://example.com/ws",
			"rpc_node.http_url":                    "https://example.com/rpc",
			"rpc_node.max_addresses_per_log_query": 500,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := rpc_node_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.Equal(t, 500, settings.MaxAddressesPerLogQuery)
		})
	})

	t.Run("Given urls with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.ws_url":   "  wss://example.com/ws  ",
			"rpc_node.http_url": "  https://example.com/rpc  ",
		}, "."), nil)

		t.Run("Should trim them", func(t *testing.T) {
			t.Parallel()

			settings := rpc_node_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.Equal(t, "wss://example.com/ws", settings.WsUrl)
			assert.Equal(t, "https://example.com/rpc", settings.HttpUrl)
		})
	})
}
