package use_cases

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetConfigInterface interface {
	Get() (*responses.ConfigResponse, error)
}

type GetConfig struct {
	blockchainConfigSettings *settings.BlockchainConfigSettings
	walletSettings           *settings.WalletSettings
	ipfsGatewaysSettings     *settings.IpfsGatewaysSettings
}

func NewGetConfig(
	blockchainConfigSettings *settings.BlockchainConfigSettings,
	walletSettings *settings.WalletSettings,
	ipfsGatewaysSettings *settings.IpfsGatewaysSettings,
) *GetConfig {
	return &GetConfig{
		blockchainConfigSettings: blockchainConfigSettings,
		walletSettings:           walletSettings,
		ipfsGatewaysSettings:     ipfsGatewaysSettings,
	}
}

func (uc *GetConfig) Get() (*responses.ConfigResponse, error) {
	evmChains := make(map[string]responses.EvmChainConfig, len(uc.blockchainConfigSettings.EvmChains))
	for chainId, chain := range uc.blockchainConfigSettings.EvmChains {
		evmChains[chainId] = responses.EvmChainConfig{
			CollectionsManagerAddress: strings.ToLower(chain.CollectionsManagerAddress),
			CollectionCreatorAddress:  strings.ToLower(chain.CollectionCreatorAddress),
			ScannerBaseUrl:            chain.ScannerBaseUrl,
		}
	}

	return &responses.ConfigResponse{
		EvmChains:               evmChains,
		OrchestratorAddress:     strings.ToLower(uc.walletSettings.Address.Hex()),
		IpfsGatewayUrlTemplates: uc.ipfsGatewaysSettings.GatewayUrlTemplates,
	}, nil
}
