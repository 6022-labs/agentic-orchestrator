package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationRepository struct {
	AbstractRepository[models.Conversation]
}

func NewConversationRepository(client db_clients.DatabaseClientInterface) *ConversationRepository {
	return &ConversationRepository{
		AbstractRepository: *NewAbstractRepository[models.Conversation](client),
	}
}

func (r *ConversationRepository) GetById(id uuid.UUID) (*models.Conversation, error) {
	var entity models.Conversation
	err := r.client.GetDB().
		Where("id = ?", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
