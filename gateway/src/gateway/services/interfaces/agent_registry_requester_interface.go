package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type ListRegistryAgentsRequest struct {
	Page            *uint
	ElementsPerPage *uint
	SortBy          *string
	SortDirection   *string
	// IsMinted filters on mint state: true = minted only, false = proposals only, nil = all.
	IsMinted *bool
}

type AgentRegistryRequesterInterface interface {
	ListAgents(request *ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error)
	GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.AgentResponse, error)
	ListAgentCollections() ([]*responses.AgentCollectionResponse, error)
}
