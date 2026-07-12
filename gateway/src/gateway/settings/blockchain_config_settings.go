package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const BlockchainConfigSettingsKey = "blockchain"

// EvmChainSettings is one EVM chain's contract deployment.
type EvmChainSettings struct {
	CollectionsManagerAddress string `koanf:"collections_manager_address"`
	CollectionCreatorAddress  string `koanf:"collection_creator_address"`
}

// BlockchainConfigSettings feeds the public GET /config, keyed by chain id.
type BlockchainConfigSettings struct {
	EvmChains map[string]EvmChainSettings `koanf:"evm_chains"`
}

func NewBlockchainConfigSettings(logger *zap.Logger, k *koanf.Koanf) *BlockchainConfigSettings {
	settings := BlockchainConfigSettings{EvmChains: map[string]EvmChainSettings{}}
	if err := k.Unmarshal(BlockchainConfigSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal blockchain config settings", zap.Error(err))
	}

	if len(settings.EvmChains) == 0 {
		logger.Warn("blockchain.evm_chains is empty, the console will fall back to its defaults")
	}

	return &settings
}
