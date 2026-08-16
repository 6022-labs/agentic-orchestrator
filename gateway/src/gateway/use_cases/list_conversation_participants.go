package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/google/uuid"
)

type ListConversationParticipantsInterface interface {
	List(conversationId uuid.UUID) ([]*responses.BrokerParticipantResponse, error)
}

type ListConversationParticipants struct {
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface
}

func NewListConversationParticipants(
	conversationBrokerRequester interfaces.ConversationBrokerRequesterInterface,
) *ListConversationParticipants {
	return &ListConversationParticipants{conversationBrokerRequester: conversationBrokerRequester}
}

func (uc *ListConversationParticipants) List(conversationId uuid.UUID) ([]*responses.BrokerParticipantResponse, error) {
	return uc.conversationBrokerRequester.ListConversationParticipants(conversationId)
}
