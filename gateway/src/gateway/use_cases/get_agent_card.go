package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
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
		return nil, errors.NewUnavailableError("agent_registry_unavailable", "agent registry is unavailable", err)
	}
	if agent == nil {
		return nil, errors.NewNotFoundError("agent_not_found", "agent not found: "+request.AgentCollectionAddress)
	}
	if agent.EnsDomain == "" {
		return nil, errors.NewNotFoundError(
			"agent_ens_domain_not_found",
			"agent has no ens domain: "+request.AgentCollectionAddress,
		)
	}

	response, err := uc.conversationBrokerRequester.GetAgentCard(ctx, agent.EnsDomain)
	if err != nil {
		return nil, errors.NewUnavailableError("agent_node_unreachable", "agent node is unreachable", err)
	}

	return response, nil
}
