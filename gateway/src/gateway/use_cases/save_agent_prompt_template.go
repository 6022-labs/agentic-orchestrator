package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type SaveAgentPromptTemplateInterface interface {
	Save(
		ctx context.Context,
		request *requests.SaveAgentPromptTemplateRequest,
	) (*responses.AgentPromptTemplateResponse, error)
}

type SaveAgentPromptTemplate struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewSaveAgentPromptTemplate(
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *SaveAgentPromptTemplate) Save(
	ctx context.Context,
	request *requests.SaveAgentPromptTemplateRequest,
) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return uc.conversationBrokerRequester.SaveAgentPromptTemplate(ctx,
		request.AgentCollectionAddress,
		request.ParsedCollectionTokenId(),
		request.Key,
		request.Content,
	)
}
