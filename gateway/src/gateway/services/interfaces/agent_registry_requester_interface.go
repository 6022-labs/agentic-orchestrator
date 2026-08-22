package interfaces

import (
	"context"

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
	ListAgents(ctx context.Context, request *ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error)
	GetAgent(
		ctx context.Context,
		agentCollectionAddress string,
		agentCollectionTokenId uint64,
	) (*responses.AgentResponse, error)
	ListAgentCollections(ctx context.Context) ([]*responses.AgentCollectionResponse, error)
}
