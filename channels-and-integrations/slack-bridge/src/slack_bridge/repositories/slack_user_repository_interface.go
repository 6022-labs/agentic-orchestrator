package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackUserRepositoryInterface interface {
	Save(entity *models.SlackUser) (*models.SlackUser, error)
	GetById(id string) (*models.SlackUser, error)
	GetByPairedHumanAgent(collectionAddress string, collectionTokenId uint64) (*models.SlackUser, error)
}
