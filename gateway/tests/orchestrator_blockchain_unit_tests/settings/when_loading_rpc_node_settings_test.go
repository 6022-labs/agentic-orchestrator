package settings_test

import (
	"testing"

	rpc_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRpcNodeSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"rpc_node.ws_url":       "wss://ws.example",
			"rpc_node.http_url":     "https://http.example",
			"rpc_node.ens_http_url": "https://ens.example",
		}, "."), nil)

		settings := rpc_settings.NewRpcNodeSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "wss://ws.example", settings.WsUrl)
		assert.Equal(t, "https://http.example", settings.HttpUrl)
		assert.Equal(t, "https://ens.example", settings.ENSHttpUrl)
	})
}
