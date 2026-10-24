package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type ForceConcludeConversationInterface interface {
	Conclude(ctx context.Context, request *requests.ConversationByIdRequest) error
}

type ForceConcludeConversation struct {
	conversationAccessVerifier  interfaces.ConversationAccessVerifierInterface
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewForceConcludeConversation(
	conversationAccessVerifier interfaces.ConversationAccessVerifierInterface,
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *ForceConcludeConversation {
	return &ForceConcludeConversation{
		conversationAccessVerifier:  conversationAccessVerifier,
		conversationBrokerRequester: conversationBrokerRequester,
	}
}

func (uc *ForceConcludeConversation) Conclude(ctx context.Context, request *requests.ConversationByIdRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	callerId, err := helpers.RequireCallerId(ctx)
	if err != nil {
		return err
	}

	if err := uc.conversationAccessVerifier.EnsureAccess(ctx,
		callerId,
		request.ConversationId,
	); err != nil {
		return err
	}

	return uc.conversationBrokerRequester.ForceConcludeConversation(ctx, request.ConversationId)
}
