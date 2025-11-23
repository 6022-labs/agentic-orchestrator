package repositories

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackMessageRepositoryInterface interface {
	Save(entity *models.SlackMessage) (*models.SlackMessage, error)
	SaveMany(entities []*models.SlackMessage) ([]*models.SlackMessage, error)
	GetConversationByMessageId(id string) ([]*models.SlackMessage, error)
	GetByThreadIdAndChannelId(threadId, channelId string) ([]*models.SlackMessage, error)
	GetProcessableMessagesForAgentAndThread(agent *models.SlackAgent, threadId, channelId string) ([]*models.SlackMessage, error)
}
