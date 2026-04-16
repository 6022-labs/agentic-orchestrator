package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackMessageLinkRepositoryInterface interface {
	Save(entity *models.SlackMessageLink) (*models.SlackMessageLink, error)
	GetBySlackIds(slackChannelId, slackMessageTs string) (*models.SlackMessageLink, error)
}
