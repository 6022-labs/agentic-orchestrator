package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationParticipantRepositoryInterface interface {
	Save(entity *models.ConversationParticipant) (*models.ConversationParticipant, error)
	ListByConversationId(conversationId uuid.UUID) ([]*models.ConversationParticipant, error)
	DeleteByConversationIdAndParticipantId(conversationId uuid.UUID, participantId uuid.UUID) error
}
