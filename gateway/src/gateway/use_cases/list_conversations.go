package use_cases

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type ListConversationsInterface interface {
	List(ctx context.Context, page, elementsPerPage *uint) (*responses.ListBrokerConversationsResponse, error)
}

type ListConversations struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewListConversations(
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *ListConversations {
	return &ListConversations{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *ListConversations) List(
	ctx context.Context,
	page, elementsPerPage *uint,
) (*responses.ListBrokerConversationsResponse, error) {
	return uc.conversationBrokerRequester.ListConversations(ctx, page, elementsPerPage)
}
