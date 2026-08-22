package repositories

import (
	"context"

	"errors"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository struct {
	AbstractRepository[models.Conversation]
}

func NewConversationRepository(client db_clients.DatabaseClientInterface) *ConversationRepository {
	return &ConversationRepository{
		AbstractRepository: *NewAbstractRepository[models.Conversation](client),
	}
}

func (r *ConversationRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	var entity models.Conversation
	err := r.client.GetDB().WithContext(ctx).
		Where("id = ?", id).
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

func (r *ConversationRepository) List(ctx context.Context, page, elementsPerPage uint) ([]*models.Conversation, error) {
	var entities []*models.Conversation
	err := r.client.GetDB().WithContext(ctx).
		Order("updated_at DESC").
		Offset(int((page - 1) * elementsPerPage)).
		Limit(int(elementsPerPage)).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *ConversationRepository) Count(ctx context.Context) (uint, error) {
	var count int64
	err := r.client.GetDB().WithContext(ctx).
		Model(&models.Conversation{}).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (r *ConversationRepository) ListProcessableForParticipant(
	ctx context.Context,
	participant *models.Participant,
	role *enums.ParticipantRole,
	facilitatorTickIntervalSeconds uint64,
) ([]*models.Conversation, error) {
	query := r.client.GetDB().WithContext(ctx).Model(&models.Conversation{}).
		Joins("INNER JOIN conversation_participants ON conversation_participants.conversation_id = conversations.id").
		Where("conversation_participants.participant_id = ?", participant.ID).
		Where("conversations.status = ?", enums.ConversationStatusActive)

	// A nil role is a user, who has no role-specific wake rule.
	if role != nil {
		ensMention := ""
		if participant.EnsDomain != nil {
			ensMention = "%" + *participant.EnsDomain + "%"
		}

		switch *role {
		case enums.ParticipantRoleFacilitator:
			// Wake only when the facilitator's own last message is older than
			// the tick interval (or they never posted).
			lastOwnMessageEpoch := r.client.GetDB().WithContext(ctx).
				Model(&models.ConversationMessage{}).
				Select("MAX(EXTRACT(EPOCH FROM conversation_messages.created_at))").
				Where("conversation_messages.conversation_id = conversations.id").
				Where("conversation_messages.participant_id = ?", participant.ID)

			query = query.Where(
				fmt.Sprintf("COALESCE((?), 0) < (EXTRACT(EPOCH FROM NOW()) - %d)", facilitatorTickIntervalSeconds),
				lastOwnMessageEpoch,
			)

		case enums.ParticipantRoleHuman:
			// Wake if any message in the conversation mentions the human's ENS
			// OR was authored by them.
			mention := r.client.GetDB().WithContext(ctx).Model(&models.ConversationMessage{}).
				Select("1").
				Where("conversation_messages.conversation_id = conversations.id").
				Where("conversation_messages.content ILIKE ?", ensMention)
			own := r.client.GetDB().WithContext(ctx).Model(&models.ConversationMessage{}).
				Select("1").
				Where("conversation_messages.conversation_id = conversations.id").
				Where("conversation_messages.participant_id = ?", participant.ID)

			query = query.Where("EXISTS (?) OR EXISTS (?)", mention, own)

		default:
			// Clone / Expert: wake only when a new ENS mention has landed AFTER
			// this agent's last own reply. Prevents infinite re-processing.
			mentionSinceLastReply := r.client.GetDB().WithContext(ctx).Model(&models.ConversationMessage{}).
				Select("1").
				Where("conversation_messages.conversation_id = conversations.id").
				Where("conversation_messages.content ILIKE ?", ensMention).
				Where(`
					conversation_messages.created_at > COALESCE(
						(
							SELECT MAX(cm2.created_at)
							FROM conversation_messages cm2
							WHERE cm2.conversation_id = conversations.id
								AND cm2.participant_id = ?
						),
						'epoch'::timestamptz
					)`,
					participant.ID,
				)

			query = query.Where("EXISTS (?)", mentionSinceLastReply)
		}
	}

	var entities []*models.Conversation
	err := query.
		Distinct(
			"conversations.id",
			"conversations.updated_at",
			"conversations.created_at",
			"conversations.bridge_id",
			"conversations.description",
			"conversations.status",
			"conversations.concluded_at",
			"conversations.unconcluded_at",
		).
		Order("conversations.updated_at ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
