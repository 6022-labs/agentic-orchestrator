package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationMessageRepository struct {
	AbstractRepository[models.ConversationMessage]
}

func NewConversationMessageRepository(client db_clients.DatabaseClientInterface) *ConversationMessageRepository {
	return &ConversationMessageRepository{
		AbstractRepository: *NewAbstractRepository[models.ConversationMessage](client),
	}
}

func (r *ConversationMessageRepository) ListByConversationId(conversationId uuid.UUID) ([]*models.ConversationMessage, error) {
	var entities []*models.ConversationMessage
	err := r.client.GetDB().
		Where("conversation_id = ?", conversationId).
		Order("created_at ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
