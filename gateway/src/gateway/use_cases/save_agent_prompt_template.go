package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
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
	agentAccessVerifier         interfaces.AgentAccessVerifierInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewSaveAgentPromptTemplate(
	agentAccessVerifier interfaces.AgentAccessVerifierInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{
		agentAccessVerifier:         agentAccessVerifier,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (uc *SaveAgentPromptTemplate) Save(
	ctx context.Context,
	request *requests.SaveAgentPromptTemplateRequest,
) (*responses.AgentPromptTemplateResponse, error) {
	if err := request.ValidateAndSanitize(); err != nil {
		return nil, err
	}
	callerId, err := helpers.RequireCallerId(ctx)
	if err != nil {
		return nil, err
	}

	if err := uc.agentAccessVerifier.EnsureAccess(ctx,
		callerId,
		request.AgentCollectionAddress,
		request.ParsedCollectionTokenId(),
	); err != nil {
		return nil, err
	}

	return uc.conversationBrokerRequester.SaveAgentPromptTemplate(ctx,
		request.AgentCollectionAddress,
		request.ParsedCollectionTokenId(),
		request.Key,
		request.Content,
	)
}
