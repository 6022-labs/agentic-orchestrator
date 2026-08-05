package settings_test

import (
	"testing"

	registry_common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	agent_collections_manager_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionsManagerSettings(t *testing.T) {
	t.Parallel()

	blockchainSettings := &registry_common_settings.BlockchainSettings{ChainId: 80002}

	t.Run("Given an explicitly configured address", func(t *testing.T) {
		t.Parallel()

		address := "0x0000000000000000000000000000000000000001"

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_collections_manager.smart_contract_address": address,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := agent_collections_manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t, common.HexToAddress(address), settings.SmartContractAddress)
		})
	})

	t.Run("Given no configured address", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to the built-in address for the chain", func(t *testing.T) {
			t.Parallel()

			settings := agent_collections_manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), k, blockchainSettings)

			assert.NotNil(t, settings)
			assert.NotEqual(t, common.Address{}, settings.SmartContractAddress)
		})
	})
}
