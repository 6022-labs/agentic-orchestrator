package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
)

type SlackUserRepository struct {
	AbstractRepository[models.SlackUser]
}

func NewSlackUserRepository(
	client *slack_bridge_db.DatabaseClient,
) *SlackUserRepository {
	return &SlackUserRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackUser](client),
	}
}

func (r *SlackUserRepository) GetById(id string) (*models.SlackUser, error) {
	var entity models.SlackUser
	err := r.client.
		Where("id = ?", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *SlackUserRepository) GetByPairedHumanAgent(collectionAddress string, collectionTokenId uint64) (*models.SlackUser, error) {
	var entity models.SlackUser
	err := r.client.
		Where("paired_human_agent_collection_address = ? AND paired_human_agent_collection_token_id = ?", collectionAddress, collectionTokenId).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
