package settings_test

import (
	"testing"

	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	token_6022_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_blockchain/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingToken6022Settings(t *testing.T) {
	t.Run("Given a chain id present in the built-in registry", func(t *testing.T) {
		blockchainSettings := &common_settings.BlockchainSettings{ChainId: 80002}

		t.Run("Should resolve the well-known address for that chain", func(t *testing.T) {
			settings := token_6022_settings.NewToken6022Settings(zap.NewNop(), blockchainSettings)

			assert.NotNil(t, settings)
			assert.Equal(t,
				common.HexToAddress("0x26adaf578361953e8F1aC86D3e4Ca153493af5cB"),
				settings.SmartContractAddress,
			)
		})
	})
}
