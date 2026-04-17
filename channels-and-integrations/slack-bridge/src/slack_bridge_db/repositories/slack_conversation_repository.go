package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/google/uuid"
)

type SlackConversationRepository struct {
	AbstractRepository[models.SlackConversation]
}

func NewSlackConversationRepository(client *slack_bridge_db.DatabaseClient) *SlackConversationRepository {
	return &SlackConversationRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackConversation](client),
	}
}

func (r *SlackConversationRepository) GetBySlackIds(slackChannelId, slackThreadTs string) (*models.SlackConversation, error) {
	var entity models.SlackConversation
	err := r.client.
		Where("slack_channel_id = ? AND slack_thread_ts = ?", slackChannelId, slackThreadTs).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SlackConversationRepository) GetByBrokerConversationId(brokerConversationId uuid.UUID) (*models.SlackConversation, error) {
	var entity models.SlackConversation
	err := r.client.
		Where("broker_conversation_id = ?", brokerConversationId).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SlackConversationRepository) ListByChannelId(slackChannelId string) ([]*models.SlackConversation, error) {
	var entities []*models.SlackConversation
	err := r.client.
		Where("slack_channel_id = ?", slackChannelId).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
