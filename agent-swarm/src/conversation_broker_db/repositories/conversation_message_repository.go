package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
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

func (r *ConversationMessageRepository) ListByConversationId(
	ctx context.Context,
	conversationId uuid.UUID,
) ([]*models.ConversationMessage, error) {
	var entities []*models.ConversationMessage
	err := r.client.GetDB().WithContext(ctx).
		Where("conversation_id = ?", conversationId).
		Order("created_at ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *ConversationMessageRepository) ListProcessableForParticipantAndConversation(
	ctx context.Context,
	participant *models.Participant,
	role *enums.ParticipantRole,
	conversationId uuid.UUID,
) ([]*models.ConversationMessage, error) {
	query := r.client.GetDB().WithContext(ctx).Model(&models.ConversationMessage{}).
		Where("conversation_id = ?", conversationId)

	// A nil role is a user, whose messages are not filtered.
	if role != nil {
		ensMention := "%" + *participant.EnsDomain + "%"

		switch *role {
		case enums.ParticipantRoleClone, enums.ParticipantRoleExpert:
			query = query.Where("content ILIKE ?", ensMention)
		case enums.ParticipantRoleHuman:
			query = query.Where("participant_id = ? OR content ILIKE ?", participant.ID, ensMention)
		}
	}

	var entities []*models.ConversationMessage
	err := query.Order("created_at ASC").Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
