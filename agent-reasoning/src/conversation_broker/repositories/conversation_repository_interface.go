package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationRepositoryInterface interface {
	Save(entity *models.Conversation) (*models.Conversation, error)
	GetById(id uuid.UUID) (*models.Conversation, error)
}
