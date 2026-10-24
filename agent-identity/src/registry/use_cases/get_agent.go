package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

type GetAgentInterface interface {
	Get(ctx context.Context, request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error)
}

type GetAgent struct {
	agentRepository        repositories.AgentRepositoryInterface
	agentAddressRepository repositories.AgentAddressRepositoryInterface
}

func NewGetAgent(
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
) *GetAgent {
	return &GetAgent{
		agentRepository:        agentRepository,
		agentAddressRepository: agentAddressRepository,
	}
}

func (s *GetAgent) Get(
	ctx context.Context,
	request *requests.AgentByCollectionTokenIdRequest,
) (*responses.AgentResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	agent, err := s.agentRepository.GetByCollectionAddressAndCollectionTokenId(ctx,
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, nil
	}

	agentAddresses, err := s.agentAddressRepository.GetAllByAgentId(ctx, agent.ID)
	if err != nil {
		return nil, err
	}

	return responses.AgentResponseFromModel(agent, agentAddresses), nil
}
