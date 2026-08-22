package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/google/uuid"
)

type GetConversationInterface interface {
	Get(ctx context.Context, conversationId uuid.UUID) (*responses.BrokerConversationResponse, error)
}

type GetConversation struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetConversation(conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface) *GetConversation {
	return &GetConversation{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *GetConversation) Get(
	ctx context.Context,
	conversationId uuid.UUID,
) (*responses.BrokerConversationResponse, error) {
	return uc.conversationBrokerRequester.GetConversation(ctx, conversationId)
}
