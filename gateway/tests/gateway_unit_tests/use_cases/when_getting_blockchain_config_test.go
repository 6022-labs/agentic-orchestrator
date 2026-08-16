package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type WhenGettingBlockchainConfigTestingSuite struct {
	sut *use_cases.GetBlockchainConfig
}

func WhenGettingBlockchainConfigBeforeEach(
	blockchainConfigSettings *settings.BlockchainConfigSettings,
) *WhenGettingBlockchainConfigTestingSuite {
	walletSettings := &settings.WalletSettings{
		Address: common.HexToAddress("0x2222222222222222222222222222222222222222"),
	}

	return &WhenGettingBlockchainConfigTestingSuite{
		sut: use_cases.NewGetBlockchainConfig(blockchainConfigSettings, walletSettings),
	}
}

func TestWhenGettingBlockchainConfig(t *testing.T) {
	t.Parallel()

	t.Run("Given chains are configured", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenGettingBlockchainConfigTestingSuite {
			return WhenGettingBlockchainConfigBeforeEach(&settings.BlockchainConfigSettings{
				EvmChains: map[string]settings.EvmChainSettings{
					"80002": {
						CollectionsManagerAddress: "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
						CollectionCreatorAddress:  "0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16",
						ScannerBaseUrl:            "https://amoy.polygonscan.com",
					},
				},
			})
		}

		t.Run("Should expose them keyed by chain id", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			resp, err := suite.sut.Get()

			assert.NoError(t, err)
			if assert.Len(t, resp.EvmChains, 1) {
				chain := resp.EvmChains["80002"]
				assert.Equal(t, "0x39DD75a025145770a8b8278Cc8BA705116443Ea7", chain.CollectionsManagerAddress)
				assert.Equal(t, "0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16", chain.CollectionCreatorAddress)
				assert.Equal(t, "https://amoy.polygonscan.com", chain.ScannerBaseUrl)
			}
		})

		t.Run("Should expose the orchestrator wallet address", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			resp, err := suite.sut.Get()

			assert.NoError(t, err)
			assert.Equal(t, "0x2222222222222222222222222222222222222222", resp.OrchestratorAddress)
		})
	})

	t.Run("Given no chain is configured", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenGettingBlockchainConfigTestingSuite {
			return WhenGettingBlockchainConfigBeforeEach(
				&settings.BlockchainConfigSettings{EvmChains: map[string]settings.EvmChainSettings{}},
			)
		}

		t.Run("Should return an empty map", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			resp, err := suite.sut.Get()

			assert.NoError(t, err)
			assert.Empty(t, resp.EvmChains)
		})
	})
}
