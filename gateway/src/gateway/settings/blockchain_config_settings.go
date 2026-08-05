package settings

import (
	_ "embed"
	"encoding/json"
	"strconv"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const BlockchainConfigSettingsKey = "blockchain"

//go:embed evm_chains.json
var evmChainsJson []byte

type evmChainEntry struct {
	ChainId                   int64  `json:"chainId"`
	Name                      string `json:"name"`
	CollectionsManagerAddress string `json:"collectionsManagerAddress"`
	CollectionCreatorAddress  string `json:"collectionCreatorAddress"`
	ScannerBaseUrl            string `json:"scannerBaseUrl"`
}

type evmChainsFile struct {
	Chains []evmChainEntry `json:"chains"`
}

// EvmChainSettings is one EVM chain's contract deployment and explorer.
type EvmChainSettings struct {
	CollectionsManagerAddress string `koanf:"collections_manager_address"`
	CollectionCreatorAddress  string `koanf:"collection_creator_address"`
	ScannerBaseUrl            string `koanf:"scanner_base_url"`
}

// BlockchainConfigSettings feeds the public GET /config, keyed by chain id.
type BlockchainConfigSettings struct {
	EvmChains map[string]EvmChainSettings `koanf:"evm_chains"`
}

func NewBlockchainConfigSettings(logger *zap.Logger, k *koanf.Koanf) *BlockchainConfigSettings {
	var file evmChainsFile
	if err := json.Unmarshal(evmChainsJson, &file); err != nil {
		logger.Fatal("failed to parse built-in evm_chains.json", zap.Error(err))
	}

	wellKnown := map[string]evmChainEntry{}
	for _, chain := range file.Chains {
		wellKnown[strconv.FormatInt(chain.ChainId, 10)] = chain
	}

	settings := BlockchainConfigSettings{EvmChains: map[string]EvmChainSettings{}}
	if err := k.Unmarshal(BlockchainConfigSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal blockchain config settings", zap.Error(err))
	}

	if len(settings.EvmChains) == 0 {
		for chainId, chain := range wellKnown {
			if len(chain.CollectionsManagerAddress) == 0 || len(chain.CollectionCreatorAddress) == 0 {
				continue
			}
			settings.EvmChains[chainId] = EvmChainSettings{}
		}
	}

	for chainId, configured := range settings.EvmChains {
		chain, ok := wellKnown[chainId]
		if !ok {
			continue
		}
		if len(configured.CollectionsManagerAddress) == 0 {
			configured.CollectionsManagerAddress = chain.CollectionsManagerAddress
		}
		if len(configured.CollectionCreatorAddress) == 0 {
			configured.CollectionCreatorAddress = chain.CollectionCreatorAddress
		}
		if len(configured.ScannerBaseUrl) == 0 {
			configured.ScannerBaseUrl = chain.ScannerBaseUrl
		}
		settings.EvmChains[chainId] = configured
	}

	if len(settings.EvmChains) == 0 {
		logger.Warn("blockchain.evm_chains is empty, the console will fall back to its defaults")
	}

	return &settings
}
