package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

type ListAgentCollectionsInterface interface {
	List(ctx context.Context) ([]*responses.AgentCollectionResponse, error)
}

type ListAgentCollections struct {
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

func NewListAgentCollections(
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
) *ListAgentCollections {
	return &ListAgentCollections{agentCollectionRepository: agentCollectionRepository}
}

func (s *ListAgentCollections) List(ctx context.Context) ([]*responses.AgentCollectionResponse, error) {
	agentCollections, err := s.agentCollectionRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*responses.AgentCollectionResponse, len(agentCollections))
	for i, c := range agentCollections {
		out[i] = &responses.AgentCollectionResponse{Address: c.Address, Name: c.Name, Symbol: c.Symbol}
	}

	return out, nil
}
