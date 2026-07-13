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
	address := "0x0000000000000000000000000000000000000001"

	t.Run("Loads settings from an explicit address", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_collections_manager.smart_contract_address": address,
		}, "."), nil)

		settings := agent_collections_manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), k, &registry_common_settings.BlockchainSettings{ChainId: 80002})

		assert.NotNil(t, settings)
		assert.True(t, common.HexToAddress(address) == settings.SmartContractAddress)
	})

	t.Run("Falls back to the built-in address when none is configured", func(t *testing.T) {
		k := koanf.New(".")

		settings := agent_collections_manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), k, &registry_common_settings.BlockchainSettings{ChainId: 80002})

		assert.NotNil(t, settings)
		assert.True(t, common.HexToAddress("0x33Fec4914A2f3e81Af2665B60b15Fb68c438e0d8") == settings.SmartContractAddress)
	})
}
