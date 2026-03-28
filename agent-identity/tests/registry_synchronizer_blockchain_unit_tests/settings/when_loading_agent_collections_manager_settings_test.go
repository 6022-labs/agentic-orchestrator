package settings_test

import (
	"testing"

	agent_collections_manager_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionsManagerSettings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"

	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_collections_manager.smart_contract_address": address,
		}, "."), nil)

		settings := agent_collections_manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.True(t, common.HexToAddress(address) == settings.SmartContractAddress)
	})
}
