package settings_test

import (
	"testing"

	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	agent_collections_manager_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionsManagerSettings(t *testing.T) {
	t.Run("Given a chain id present in the built-in registry", func(t *testing.T) {
		blockchainSettings := &common_settings.BlockchainSettings{ChainId: 80002}

		t.Run("Should resolve the well-known address for that chain", func(t *testing.T) {
			settings := agent_collections_manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t,
				common.HexToAddress("0x39DD75a025145770a8b8278Cc8BA705116443Ea7"),
				settings.SmartContractAddress,
			)
		})
	})
}
