package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/responses"
)

type ListAgentCollectionsInterface interface {
	List() ([]*responses.AgentCollectionResponse, error)
}

type ListAgentCollections struct {
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface
}

func NewListAgentCollections(agentCollectionRepository repositories.AgentCollectionRepositoryInterface) *ListAgentCollections {
	return &ListAgentCollections{
		agentCollectionRepository: agentCollectionRepository,
	}
}

func (s *ListAgentCollections) List() ([]*responses.AgentCollectionResponse, error) {
	agentCollections, err := s.agentCollectionRepository.GetAll()
	if err != nil {
		return nil, err
	}

	agentCollectionsResponse := make([]*responses.AgentCollectionResponse, len(agentCollections))

	for index, agentCollection := range agentCollections {
		agentCollectionsResponse[index] = &responses.AgentCollectionResponse{
			Address: agentCollection.Address,
			Name:    agentCollection.Name,
			Symbol:  agentCollection.Symbol,
		}
	}

	return agentCollectionsResponse, nil
}
