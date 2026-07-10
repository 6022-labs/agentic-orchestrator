package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingEnsResolverSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"ens_resolver.rpc_url":                         "https://sepolia.example",
			"ens_resolver.registry_address":                "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
			"ens_resolver.gateway_request_timeout_seconds": uint64(20),
		}, "."), nil)

		settings := blockchain_settings.NewEnsResolverSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://sepolia.example", settings.RpcUrl)
		assert.Equal(t, uint64(20), settings.GatewayRequestTimeoutSeconds)
		assert.Equal(t, "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e", settings.RegistryCommonAddress().Hex())
	})

	t.Run("Defaults the gateway timeout and the registry address", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"ens_resolver.rpc_url": "https://sepolia.example",
		}, "."), nil)

		settings := blockchain_settings.NewEnsResolverSettings(zap.NewNop(), k)

		assert.Equal(t, uint64(10), settings.GatewayRequestTimeoutSeconds)
		assert.Equal(t, "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e", settings.RegistryCommonAddress().Hex())
	})
}
