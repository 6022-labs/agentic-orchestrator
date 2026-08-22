package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/google/uuid"
)

type ForceConcludeConversationInterface interface {
	Conclude(ctx context.Context, conversationId uuid.UUID) error
}

type ForceConcludeConversation struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewForceConcludeConversation(
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *ForceConcludeConversation {
	return &ForceConcludeConversation{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *ForceConcludeConversation) Conclude(ctx context.Context, conversationId uuid.UUID) error {
	return uc.conversationBrokerRequester.ForceConcludeConversation(ctx, conversationId)
}
