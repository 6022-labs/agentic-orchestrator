package use_cases

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentCardInterface interface {
	Get(ctx context.Context, request *requests.AgentByCollectionTokenIdRequest) (*responses.AgentCardResponse, error)
}

type GetAgentCard struct {
	agentRegistryRequester      interfaces.AgentRegistryRequesterInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetAgentCard(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *GetAgentCard {
	return &GetAgentCard{
		agentRegistryRequester:      agentRegistryRequester,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (uc *GetAgentCard) Get(
	ctx context.Context,
	request *requests.AgentByCollectionTokenIdRequest,
) (*responses.AgentCardResponse, error) {
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

	return uc.conversationBrokerRequester.GetAgentCard(ctx, agent.EnsDomain)
}
