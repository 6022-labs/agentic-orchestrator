package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
)

type IsInternalAgentInterface interface {
	IsInternal(request *requests.AgentByCollectionTokenIdRequest) (*responses.IsInternalAgentResponse, error)
}

type IsInternalAgent struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewIsInternalAgent(agentRepository repositories.AgentRepositoryInterface) *IsInternalAgent {
	return &IsInternalAgent{agentRepository: agentRepository}
}

func (u *IsInternalAgent) IsInternal(request *requests.AgentByCollectionTokenIdRequest) (*responses.IsInternalAgentResponse, error) {
	agent, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return &responses.IsInternalAgentResponse{IsInternalAgent: false}, nil
	}

	return &responses.IsInternalAgentResponse{IsInternalAgent: agent.IsInternal}, nil
}
