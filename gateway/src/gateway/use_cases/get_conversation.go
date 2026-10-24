package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetConversationInterface interface {
	Get(ctx context.Context, request *requests.ConversationByIdRequest) (*responses.BrokerConversationResponse, error)
}

type GetConversation struct {
	conversationAccessVerifier  interfaces.ConversationAccessVerifierInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetConversation(
	conversationAccessVerifier interfaces.ConversationAccessVerifierInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *GetConversation {
	return &GetConversation{
		conversationAccessVerifier:  conversationAccessVerifier,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (uc *GetConversation) Get(
	ctx context.Context,
	request *requests.ConversationByIdRequest,
) (*responses.BrokerConversationResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	callerId, err := helpers.RequireCallerId(ctx)
	if err != nil {
		return nil, err
	}

	if err := uc.conversationAccessVerifier.EnsureAccess(ctx,
		callerId,
		request.ConversationId,
	); err != nil {
		return nil, err
	}

	return uc.conversationBrokerRequester.GetConversation(ctx, request.ConversationId)
}
