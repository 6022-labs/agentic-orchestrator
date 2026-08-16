package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/google/uuid"
)

type GetConversationMessagesInterface interface {
	Get(conversationId uuid.UUID) ([]*responses.BrokerMessageResponse, error)
}

type GetConversationMessages struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewGetConversationMessages(
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *GetConversationMessages {
	return &GetConversationMessages{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *GetConversationMessages) Get(conversationId uuid.UUID) ([]*responses.BrokerMessageResponse, error) {
	return uc.conversationBrokerRequester.GetConversationMessages(conversationId)
}
