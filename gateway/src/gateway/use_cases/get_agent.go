package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentInterface interface {
	Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error)
}

type GetAgent struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewGetAgent(agentRegistryRequester interfaces.AgentRegistryRequesterInterface) *GetAgent {
	return &GetAgent{agentRegistryRequester: agentRegistryRequester}
}

func (s *GetAgent) Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error) {
	return s.agentRegistryRequester.GetAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
}
