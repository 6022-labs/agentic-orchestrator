package use_cases

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"gorm.io/gorm"
)

type GetMintedAgentByEnsDomainInterface interface {
	Get(request *requests.GetAgentByEnsDomainRequest) (*responses.AgentResponse, error)
}

type GetMintedAgentByEnsDomain struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewGetMintedAgentByEnsDomain(agentRepository repositories.AgentRepositoryInterface) *GetMintedAgentByEnsDomain {
	return &GetMintedAgentByEnsDomain{agentRepository: agentRepository}
}

func (s *GetMintedAgentByEnsDomain) Get(request *requests.GetAgentByEnsDomainRequest) (*responses.AgentResponse, error) {
	agent, err := s.agentRepository.GetMintedByEnsDomain(request.EnsDomain)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if agent == nil {
		return nil, nil
	}

	return responses.AgentResponseFromModel(agent), nil
}
