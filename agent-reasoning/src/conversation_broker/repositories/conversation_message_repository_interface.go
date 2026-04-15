package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationMessageRepositoryInterface interface {
	Save(entity *models.ConversationMessage) (*models.ConversationMessage, error)
	ListByConversationId(conversationId uuid.UUID) ([]*models.ConversationMessage, error)
}
