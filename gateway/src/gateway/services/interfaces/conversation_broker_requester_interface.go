package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/google/uuid"
)

// ConversationBrokerRequesterInterface is the gateway's read-side port to the
// conversation broker (conversations, prompt templates, agent cards).
type ConversationBrokerRequesterInterface interface {
	ListConversations(
		ctx context.Context,
		page, elementsPerPage *uint,
	) (*responses.ListBrokerConversationsResponse, error)
	GetConversation(ctx context.Context, conversationId uuid.UUID) (*responses.BrokerConversationResponse, error)
	GetConversationMessages(ctx context.Context, conversationId uuid.UUID) ([]*responses.BrokerMessageResponse, error)
	ListConversationParticipants(
		ctx context.Context,
		conversationId uuid.UUID,
	) ([]*responses.BrokerParticipantResponse, error)
	ForceConcludeConversation(ctx context.Context, conversationId uuid.UUID) error
	GetAgentPromptTemplate(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
		key string,
	) (*responses.AgentPromptTemplateResponse, error)
	SaveAgentPromptTemplate(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
		key, content string,
	) (*responses.AgentPromptTemplateResponse, error)
	GetAgentCard(ctx context.Context, ensDomain string) (*responses.AgentCardResponse, error)
	GetAgentUrl(ctx context.Context, ensDomain string) (*responses.AgentUrlResponse, error)
}
