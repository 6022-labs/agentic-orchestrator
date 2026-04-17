package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackChannelRepositoryInterface interface {
	Save(entity *models.SlackChannel) (*models.SlackChannel, error)
	GetById(id string) (*models.SlackChannel, error)
	AddAgent(entity *models.SlackChannel, agent *models.SlackAgent) (*models.SlackChannel, error)
	RemoveAgent(entity *models.SlackChannel, agent *models.SlackAgent) (*models.SlackChannel, error)
	AddUser(channelId string, userId string) error
	UpdateId(oldId string, newId string) error
}
