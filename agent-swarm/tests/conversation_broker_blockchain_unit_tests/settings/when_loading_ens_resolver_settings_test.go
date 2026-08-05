package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingEnsResolverSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ens_resolver.rpc_url":                         "https://sepolia.example",
			"ens_resolver.registry_address":                "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
			"ens_resolver.gateway_request_timeout_seconds": uint64(20),
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewEnsResolverSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://sepolia.example", settings.RpcUrl)
			assert.Equal(t, uint64(20), settings.GatewayRequestTimeoutSeconds)
			assert.Equal(t, "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e", settings.RegistryCommonAddress().Hex())
		})
	})

	t.Run("Given only the rpc url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ens_resolver.rpc_url": "https://sepolia.example",
		}, "."), nil)

		t.Run("Should fall back to the default gateway timeout and canonical registry address", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewEnsResolverSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(10), settings.GatewayRequestTimeoutSeconds)
			assert.Equal(t, "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e", settings.RegistryCommonAddress().Hex())
		})
	})
}
