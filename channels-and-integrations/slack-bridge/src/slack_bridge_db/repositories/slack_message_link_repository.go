package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
)

type SlackMessageLinkRepository struct {
	AbstractRepository[models.SlackMessageLink]
}

func NewSlackMessageLinkRepository(client *slack_bridge_db.DatabaseClient) *SlackMessageLinkRepository {
	return &SlackMessageLinkRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackMessageLink](client),
	}
}

func (r *SlackMessageLinkRepository) GetBySlackIds(slackChannelId, slackMessageTs string) (*models.SlackMessageLink, error) {
	var entity models.SlackMessageLink
	err := r.client.
		Where("slack_channel_id = ? AND slack_message_ts = ?", slackChannelId, slackMessageTs).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
