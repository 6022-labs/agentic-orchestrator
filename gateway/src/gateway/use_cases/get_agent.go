package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentInterface interface {
	Get(ctx context.Context, request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentResponse, error)
}

type GetAgent struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
}

func NewGetAgent(agentRegistryRequester interfaces.AgentRegistryRequesterInterface) *GetAgent {
	return &GetAgent{agentRegistryRequester: agentRegistryRequester}
}

func (s *GetAgent) Get(
	ctx context.Context,
	request *requests.AgentByCollectionTokenIdRequest,
) (*responses.AgentResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	return s.agentRegistryRequester.GetAgent(ctx,
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
}
