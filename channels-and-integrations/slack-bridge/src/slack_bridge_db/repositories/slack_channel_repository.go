package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
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

func (r *SlackChannelRepository) GetById(id string) (*models.SlackChannel, error) {
	var entity models.SlackChannel
	err := r.client.
		Where("id = ?", id).
		Preload("Agents").
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *SlackChannelRepository) AddAgent(entity *models.SlackChannel, agent *models.SlackAgent) (*models.SlackChannel, error) {
	r.client.GetDB().Model(entity).Association("Agents").Append(agent)
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *SlackChannelRepository) RemoveAgent(entity *models.SlackChannel, agent *models.SlackAgent) (*models.SlackChannel, error) {
	r.client.GetDB().Model(entity).Association("Agents").Delete(agent)
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *SlackChannelRepository) UpdateId(oldId string, newId string) error {
	return r.client.GetDB().
		Model(&models.SlackChannel{}).
		Where("id = ?", oldId).
		Update("id", newId).
		Error
}
