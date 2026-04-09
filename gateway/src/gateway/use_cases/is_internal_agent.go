package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type IsInternalAgentInterface interface {
	IsInternal(request *requests.AgentByCollectionTokenIdRequest) (*responses.IsInternalAgentResponse, error)
}

type IsInternalAgent struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewIsInternalAgent(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
) *IsInternalAgent {
	return &IsInternalAgent{
		agentRegistryRequester: agentRegistryRequester,
	}
}

func (useCase *IsInternalAgent) IsInternal(request *requests.AgentByCollectionTokenIdRequest) (*responses.IsInternalAgentResponse, error) {
	response, err := useCase.agentRegistryRequester.IsInternalAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
