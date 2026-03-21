package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

type GetAgentInterface interface {
	Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error)
}

type GetAgent struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewGetAgent(agentRepository repositories.AgentRepositoryInterface) *GetAgent {
	return &GetAgent{agentRepository: agentRepository}
}

func (s *GetAgent) Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error) {
	agent, err := s.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, nil
	}

	return responses.AgentResponseFromModel(agent), nil
}
