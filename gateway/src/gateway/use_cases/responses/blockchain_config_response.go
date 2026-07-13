package responses

// EvmChainConfig is one EVM chain's config surfaced to the console.
type EvmChainConfig struct {
	CollectionsManagerAddress string `json:"collectionsManagerAddress"`
	CollectionCreatorAddress  string `json:"collectionCreatorAddress"`
	ScannerBaseUrl            string `json:"scannerBaseUrl"`
}

// BlockchainConfigResponse is the public EVM chain config (keyed by chain id)
// the console boots from.
type BlockchainConfigResponse struct {
	EvmChains map[string]EvmChainConfig `json:"evmChains"`
}
