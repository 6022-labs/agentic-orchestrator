package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationParticipantRepository struct {
	AbstractRepository[models.ConversationParticipant]
}

func NewConversationParticipantRepository(client db_clients.DatabaseClientInterface) *ConversationParticipantRepository {
	return &ConversationParticipantRepository{
		AbstractRepository: *NewAbstractRepository[models.ConversationParticipant](client),
	}
}

func (r *ConversationParticipantRepository) ListByConversationId(conversationId uuid.UUID) ([]*models.ConversationParticipant, error) {
	var entities []*models.ConversationParticipant
	err := r.client.GetDB().
		Where("conversation_id = ?", conversationId).
		Order("joined_at ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *ConversationParticipantRepository) DeleteByConversationIdAndParticipantId(
	conversationId uuid.UUID,
	participantId uuid.UUID,
) error {
	return r.client.GetDB().
		Where("conversation_id = ? AND participant_id = ?", conversationId, participantId).
		Delete(&models.ConversationParticipant{}).
		Error
}
