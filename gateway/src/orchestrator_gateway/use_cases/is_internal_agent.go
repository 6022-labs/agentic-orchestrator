package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
)

type IsInternalAgentInterface interface {
	IsInternal(request *requests.AgentByCollectionTokenIdRequest) (*responses.IsInternalAgentResponse, error)
}

type IsInternalAgent struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewIsInternalAgent(
	agentRepository repositories.AgentRepositoryInterface,
) *IsInternalAgent {
	return &IsInternalAgent{
		agentRepository: agentRepository,
	}
}

func (useCase *IsInternalAgent) IsInternal(request *requests.AgentByCollectionTokenIdRequest) (*responses.IsInternalAgentResponse, error) {
	agent, err := useCase.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return &responses.IsInternalAgentResponse{
			IsInternalAgent: false,
		}, nil
	}

	return &responses.IsInternalAgentResponse{
		IsInternalAgent: agent.IsInternal,
	}, nil
}
