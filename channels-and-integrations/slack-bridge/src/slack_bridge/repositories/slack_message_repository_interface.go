package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackMessageRepositoryInterface interface {
	Save(entity *models.SlackMessage) (*models.SlackMessage, error)
	GetBySlackIds(slackChannelId, slackMessageTs string) (*models.SlackMessage, error)
}
