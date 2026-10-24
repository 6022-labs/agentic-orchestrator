package use_cases

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentUrlInterface interface {
	Get(ctx context.Context, request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentUrlResponse, error)
}

type GetAgentUrl struct {
	agentRegistryRequester      interfaces.AgentRegistryRequesterInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetAgentUrl(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *GetAgentUrl {
	return &GetAgentUrl{
		agentRegistryRequester:      agentRegistryRequester,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (uc *GetAgentUrl) Get(
	ctx context.Context,
	request *requests.AgentByCollectionTokenIdRequest,
) (*responses.AgentUrlResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	agent, err := uc.agentRegistryRequester.GetAgent(
		ctx,
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent: %w", err)
	}
	if agent == nil {
		return nil, errors.NewNotFoundError("agent not found: " + request.AgentCollectionAddress)
	}
	if agent.EnsDomain == "" {
		return nil, errors.NewNotFoundError("agent has no ens domain: " + request.AgentCollectionAddress)
	}

	return uc.conversationBrokerRequester.GetAgentUrl(ctx, agent.EnsDomain)
}
