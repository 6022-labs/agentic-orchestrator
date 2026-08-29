package responses

// EvmChainConfig is one EVM chain's config surfaced to the console.
type EvmChainConfig struct {
	CollectionsManagerAddress string `json:"collectionsManagerAddress"`
	CollectionCreatorAddress  string `json:"collectionCreatorAddress"`
	ScannerBaseUrl            string `json:"scannerBaseUrl"`
}

// ConfigResponse is the public config the console boots from.
type ConfigResponse struct {
	EvmChains               map[string]EvmChainConfig `json:"evmChains"`
	OrchestratorAddress     string                    `json:"orchestratorAddress"`
	IpfsGatewayUrlTemplates []string                  `json:"ipfsGatewayUrlTemplates"`
}
