package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackThreadRepositoryInterface interface {
	Save(entity *models.SlackThread) (*models.SlackThread, error)
	GetByIdAndChannelId(id string, channelId string) (*models.SlackThread, error)
	GetProcessableThreadsForAgent(agent *models.SlackAgent, facilitatorTickInterval uint64) ([]*models.SlackThread, error)
	ConcludeThreadsByChannelId(channelId string) error
}
