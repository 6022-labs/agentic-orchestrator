package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentEnsRegistrySettings(t *testing.T) {
	t.Parallel()

	blockchainSettings := &blockchain_settings.BlockchainSettings{ChainId: 80002}

	t.Run("Given an explicitly configured address", func(t *testing.T) {
		t.Parallel()

		address := "0x0000000000000000000000000000000000000001"

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_ens_registry.smart_contract_address": address,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewAgentEnsRegistrySettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t, common.HexToAddress(address), settings.SmartContractAddress)
		})
	})

	t.Run("Given no configured address", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to the built-in address for the chain", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewAgentEnsRegistrySettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.NotEqual(t, common.Address{}, settings.SmartContractAddress)
		})

		t.Run("Should pick the built-in address by chain id", func(t *testing.T) {
			t.Parallel()

			amoy := blockchain_settings.NewAgentEnsRegistrySettings(zap.NewNop(), k, blockchainSettings)
			polygon := blockchain_settings.NewAgentEnsRegistrySettings(
				zap.NewNop(),
				k,
				&blockchain_settings.BlockchainSettings{ChainId: 137},
			)

			assert.NotEqual(t, amoy.SmartContractAddress, polygon.SmartContractAddress)
		})
	})
}
