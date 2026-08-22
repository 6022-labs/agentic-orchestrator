package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationParticipantRepositoryInterface interface {
	Save(ctx context.Context, entity *models.ConversationParticipant) (*models.ConversationParticipant, error)
	ListByConversationId(ctx context.Context, conversationId uuid.UUID) ([]*models.ConversationParticipant, error)
	DeleteByConversationIdAndParticipantId(ctx context.Context, conversationId uuid.UUID, participantId uuid.UUID) error
}
