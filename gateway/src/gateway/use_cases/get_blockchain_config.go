package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetBlockchainConfigInterface interface {
	Get() (*responses.BlockchainConfigResponse, error)
}

type GetBlockchainConfig struct {
	blockchainConfigSettings *settings.BlockchainConfigSettings
}

func NewGetBlockchainConfig(blockchainConfigSettings *settings.BlockchainConfigSettings) *GetBlockchainConfig {
	return &GetBlockchainConfig{blockchainConfigSettings: blockchainConfigSettings}
}

func (uc *GetBlockchainConfig) Get() (*responses.BlockchainConfigResponse, error) {
	evmChains := make(map[string]responses.EvmChainConfig, len(uc.blockchainConfigSettings.EvmChains))
	for chainId, chain := range uc.blockchainConfigSettings.EvmChains {
		evmChains[chainId] = responses.EvmChainConfig{
			CollectionsManagerAddress: chain.CollectionsManagerAddress,
			CollectionCreatorAddress:  chain.CollectionCreatorAddress,
		}
	}

	return &responses.BlockchainConfigResponse{EvmChains: evmChains}, nil
}
