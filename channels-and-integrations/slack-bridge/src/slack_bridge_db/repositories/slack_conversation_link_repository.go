package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"github.com/google/uuid"
)

type SlackConversationLinkRepository struct {
	AbstractRepository[models.SlackConversationLink]
}

func NewSlackConversationLinkRepository(client *slack_bridge_db.DatabaseClient) *SlackConversationLinkRepository {
	return &SlackConversationLinkRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackConversationLink](client),
	}
}

func (r *SlackConversationLinkRepository) GetBySlackIds(slackChannelId, slackThreadTs string) (*models.SlackConversationLink, error) {
	var entity models.SlackConversationLink
	err := r.client.
		Where("slack_channel_id = ? AND slack_thread_ts = ?", slackChannelId, slackThreadTs).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SlackConversationLinkRepository) GetByBrokerConversationId(brokerConversationId uuid.UUID) (*models.SlackConversationLink, error) {
	var entity models.SlackConversationLink
	err := r.client.
		Where("broker_conversation_id = ?", brokerConversationId).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
