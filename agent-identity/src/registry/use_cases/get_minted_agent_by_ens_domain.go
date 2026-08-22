package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

type GetMintedAgentByEnsDomainInterface interface {
	Get(ctx context.Context, request *requests.GetAgentByEnsDomainRequest) (*responses.AgentResponse, error)
}

type GetMintedAgentByEnsDomain struct {
	agentRepository        repositories.AgentRepositoryInterface
	agentAddressRepository repositories.AgentAddressRepositoryInterface
}

func NewGetMintedAgentByEnsDomain(
	agentRepository repositories.AgentRepositoryInterface,
	agentAddressRepository repositories.AgentAddressRepositoryInterface,
) *GetMintedAgentByEnsDomain {
	return &GetMintedAgentByEnsDomain{
		agentRepository:        agentRepository,
		agentAddressRepository: agentAddressRepository,
	}
}

func (s *GetMintedAgentByEnsDomain) Get(
	ctx context.Context,
	request *requests.GetAgentByEnsDomainRequest,
) (*responses.AgentResponse, error) {
	agent, err := s.agentRepository.GetMintedByEnsDomain(ctx, request.EnsDomain)
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
