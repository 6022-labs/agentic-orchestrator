package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type ListConversationParticipantsInterface interface {
	List(ctx context.Context, request *requests.ConversationByIdRequest) ([]*responses.BrokerParticipantResponse, error)
}

type ListConversationParticipants struct {
	conversationAccessVerifier  interfaces.ConversationAccessVerifierInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewListConversationParticipants(
	conversationAccessVerifier interfaces.ConversationAccessVerifierInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *ListConversationParticipants {
	return &ListConversationParticipants{
		conversationAccessVerifier:  conversationAccessVerifier,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (uc *ListConversationParticipants) List(
	ctx context.Context,
	request *requests.ConversationByIdRequest,
) ([]*responses.BrokerParticipantResponse, error) {
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

	return uc.conversationBrokerRequester.ListConversationParticipants(ctx, request.ConversationId)
}
