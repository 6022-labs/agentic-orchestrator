package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationMemoryBindingRepositoryInterface interface {
	// GetByConversationAndParticipant returns nil when the pair is still unbound.
	GetByConversationAndParticipant(
		ctx context.Context,
		conversationId uuid.UUID,
		participantId uuid.UUID,
	) (*models.ConversationMemoryBinding, error)
	// Bind records the pair's session; only a derived session is ever moved.
	Bind(ctx context.Context, binding *models.ConversationMemoryBinding) error
}
