package repositories

import (
	"context"

	"errors"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"gorm.io/gorm"
)

type SlackMessageRepository struct {
	AbstractRepository[models.SlackMessage]
}

func NewSlackMessageRepository(client *slack_bridge_db.DatabaseClient) *SlackMessageRepository {
	return &SlackMessageRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackMessage](client),
	}
}

func (r *SlackMessageRepository) GetBySlackIds(
	ctx context.Context,
	slackChannelId, slackMessageTs string,
) (*models.SlackMessage, error) {
	var entity models.SlackMessage
	err := r.client.WithContext(ctx).
		Where("slack_channel_id = ? AND slack_message_ts = ?", slackChannelId, slackMessageTs).
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
