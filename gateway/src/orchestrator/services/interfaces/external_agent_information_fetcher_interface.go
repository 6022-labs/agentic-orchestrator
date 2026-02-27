package interfaces

import "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"

type AgentMetadataInformation struct {
	LLMModel    string
	LLMPlatform enums.LLMConnectionPlatform
}

type ImportAgentInformation struct {
	PrivateKey string
}

type ExternalAgentInformationFetcherInterface interface {
	IsInternalAgent(domain string, collectionTokenId uint64) (*bool, error)
	ImportAgent(domain string, collectionTokenId uint64) (*ImportAgentInformation, error)
	FetchAgentMetadataInformation(domain string, collectionTokenId uint64) (*AgentMetadataInformation, error)
}
