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

	t.Run("Given configured urls", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"rpc_node.ws_url":   "wss://example.com/ws",
			"rpc_node.http_url": "https://example.com/rpc",
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := rpc_node_settings.NewRpcNodeSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "wss://example.com/ws", settings.WsUrl)
			assert.Equal(t, "https://example.com/rpc", settings.HttpUrl)
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
