package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type ListAgentCollectionsInterface interface {
	List(ctx context.Context) ([]*responses.AgentCollectionResponse, error)
}

type ListAgentCollections struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewListAgentCollections(agentRegistryRequester interfaces.AgentRegistryRequesterInterface) *ListAgentCollections {
	return &ListAgentCollections{
		agentRegistryRequester: agentRegistryRequester,
	}
}

func (s *ListAgentCollections) List(ctx context.Context) ([]*responses.AgentCollectionResponse, error) {
	agentCollections, err := s.agentRegistryRequester.ListAgentCollections(ctx)
	if err != nil {
		return nil, err
	}
	return agentCollections, nil
}
