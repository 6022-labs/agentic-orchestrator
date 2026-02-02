package settings_test

import (
	"testing"

	rpc_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_blockchain/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRpcNodeSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("rpc_node", map[string]any{
			"ws_url":       "wss://ws.example",
			"http_url":     "https://http.example",
			"ens_http_url": "https://ens.example",
		})

		settings := rpc_settings.NewRpcNodeSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "wss://ws.example", settings.WsUrl)
		assert.Equal(t, "https://http.example", settings.HttpUrl)
		assert.Equal(t, "https://ens.example", settings.ENSHttpUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("RPC_NODE__WS_URL", "wss://env.ws")
		t.Setenv("RPC_NODE__HTTP_URL", "https://env.http")
		t.Setenv("RPC_NODE__ENS_HTTP_URL", "https://env.ens")

		v := viper.New()
		settings := rpc_settings.NewRpcNodeSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "wss://env.ws", settings.WsUrl)
		assert.Equal(t, "https://env.http", settings.HttpUrl)
		assert.Equal(t, "https://env.ens", settings.ENSHttpUrl)
	})
}
