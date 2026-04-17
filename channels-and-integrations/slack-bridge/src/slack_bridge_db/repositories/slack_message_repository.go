package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
)

type SlackMessageRepository struct {
	AbstractRepository[models.SlackMessage]
}

func NewSlackMessageRepository(client *slack_bridge_db.DatabaseClient) *SlackMessageRepository {
	return &SlackMessageRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackMessage](client),
	}
}

func (r *SlackMessageRepository) GetBySlackIds(slackChannelId, slackMessageTs string) (*models.SlackMessage, error) {
	var entity models.SlackMessage
	err := r.client.
		Where("slack_channel_id = ? AND slack_message_ts = ?", slackChannelId, slackMessageTs).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
