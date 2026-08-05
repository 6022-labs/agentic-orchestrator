package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/google/uuid"
)

// ConversationBrokerRequesterInterface is the gateway's read-side port to the
// conversation broker (conversations, prompt templates, agent cards).
type ConversationBrokerRequesterInterface interface {
	ListConversations(page, elementsPerPage *uint) (*responses.ListBrokerConversationsResponse, error)
	GetConversation(conversationId uuid.UUID) (*responses.BrokerConversationResponse, error)
	GetConversationMessages(conversationId uuid.UUID) ([]*responses.BrokerMessageResponse, error)
	ListConversationParticipants(conversationId uuid.UUID) ([]*responses.BrokerParticipantResponse, error)
	ForceConcludeConversation(conversationId uuid.UUID) error
	GetAgentPromptTemplate(collectionAddress string, collectionTokenId uint64, key string) (*responses.AgentPromptTemplateResponse, error)
	SaveAgentPromptTemplate(collectionAddress string, collectionTokenId uint64, key, content string) (*responses.AgentPromptTemplateResponse, error)
	GetAgentCard(ensDomain string) (*responses.AgentCardResponse, error)
	GetAgentUrl(ensDomain string) (*responses.AgentUrlResponse, error)
}
