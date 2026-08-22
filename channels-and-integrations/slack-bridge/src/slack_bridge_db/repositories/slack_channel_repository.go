package repositories

import (
	"context"

	"errors"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
	"gorm.io/gorm"
)

type SlackChannelRepository struct {
	AbstractRepository[models.SlackChannel]
}

func NewSlackChannelRepository(
	client *slack_bridge_db.DatabaseClient,
) *SlackChannelRepository {
	return &SlackChannelRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackChannel](client),
	}
}

func (r *SlackChannelRepository) GetById(ctx context.Context, id string) (*models.SlackChannel, error) {
	var entity models.SlackChannel
	err := r.client.WithContext(ctx).
		Where("id = ?", id).
		Preload("Agents").
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

func (r *SlackChannelRepository) AddAgent(
	ctx context.Context,
	entity *models.SlackChannel,
	agent *models.SlackAgent,
) (*models.SlackChannel, error) {
	r.client.GetDB().WithContext(ctx).Model(entity).Association("Agents").Append(agent)
	if err := r.client.GetDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *SlackChannelRepository) RemoveAgent(
	ctx context.Context,
	entity *models.SlackChannel,
	agent *models.SlackAgent,
) (*models.SlackChannel, error) {
	r.client.GetDB().WithContext(ctx).Model(entity).Association("Agents").Delete(agent)
	if err := r.client.GetDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

// AddUser upserts the (channel, user) membership row in the slack_channel_users
// join table. Safe to call concurrently from goroutines — ON CONFLICT DO NOTHING
// makes it idempotent.
func (r *SlackChannelRepository) AddUser(ctx context.Context, channelId string, userId string) error {
	return r.client.GetDB().WithContext(ctx).
		Exec(
			`INSERT INTO slack_channel_users (slack_channel_id, slack_user_id) VALUES (?, ?) ON CONFLICT DO NOTHING`,
			channelId, userId,
		).Error
}

func (r *SlackChannelRepository) UpdateId(ctx context.Context, oldId string, newId string) error {
	return r.client.GetDB().WithContext(ctx).
		Model(&models.SlackChannel{}).
		Where("id = ?", oldId).
		Update("id", newId).
		Error
}
