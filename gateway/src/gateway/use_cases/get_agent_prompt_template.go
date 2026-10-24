package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentPromptTemplateInterface interface {
	Get(
		ctx context.Context,
		request *requests.GetAgentPromptTemplateRequest,
	) (*responses.AgentPromptTemplateResponse, error)
}

type GetAgentPromptTemplate struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetAgentPromptTemplate(
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *GetAgentPromptTemplate {
	return &GetAgentPromptTemplate{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *GetAgentPromptTemplate) Get(
	ctx context.Context,
	request *requests.GetAgentPromptTemplateRequest,
) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}

	return uc.conversationBrokerRequester.GetAgentPromptTemplate(ctx,
		request.AgentCollectionAddress,
		request.ParsedCollectionTokenId(),
		request.Key,
	)
}
