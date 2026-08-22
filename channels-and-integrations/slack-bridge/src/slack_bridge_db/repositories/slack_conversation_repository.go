package repositories

import (
	"context"

	"errors"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlackConversationRepository struct {
	AbstractRepository[models.SlackConversation]
}

func NewSlackConversationRepository(client *slack_bridge_db.DatabaseClient) *SlackConversationRepository {
	return &SlackConversationRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackConversation](client),
	}
}

func (r *SlackConversationRepository) GetBySlackIds(
	ctx context.Context,
	slackChannelId, slackThreadTs string,
) (*models.SlackConversation, error) {
	var entity models.SlackConversation
	err := r.client.WithContext(ctx).
		Where("slack_channel_id = ? AND slack_thread_ts = ?", slackChannelId, slackThreadTs).
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

func (r *SlackConversationRepository) GetByBrokerConversationId(
	ctx context.Context,
	brokerConversationId uuid.UUID,
) (*models.SlackConversation, error) {
	var entity models.SlackConversation
	err := r.client.WithContext(ctx).
		Where("broker_conversation_id = ?", brokerConversationId).
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

func (r *SlackConversationRepository) ListByChannelId(
	ctx context.Context,
	slackChannelId string,
) ([]*models.SlackConversation, error) {
	var entities []*models.SlackConversation
	err := r.client.WithContext(ctx).
		Where("slack_channel_id = ?", slackChannelId).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
