package repositories

import (
	"context"
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConversationMemoryBindingRepository struct {
	AbstractRepository[models.ConversationMemoryBinding]
}

func NewConversationMemoryBindingRepository(
	client db_clients.DatabaseClientInterface,
) *ConversationMemoryBindingRepository {
	return &ConversationMemoryBindingRepository{
		AbstractRepository: *NewAbstractRepository[models.ConversationMemoryBinding](client),
	}
}

func (r *ConversationMemoryBindingRepository) GetByConversationAndParticipant(
	ctx context.Context,
	conversationId uuid.UUID,
	participantId uuid.UUID,
) (*models.ConversationMemoryBinding, error) {
	var entity models.ConversationMemoryBinding
	err := r.client.GetDB().WithContext(ctx).
		Where("conversation_id = ? AND participant_id = ?", conversationId, participantId).
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ConversationMemoryBindingRepository) Bind(
	ctx context.Context,
	binding *models.ConversationMemoryBinding,
) error {
	return r.client.GetDB().WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "conversation_id"}, {Name: "participant_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"memory_session_id", "bound_at"}),
		}).
		Create(binding).
		Error
}
