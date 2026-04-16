package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/google/uuid"
)

type SlackConversationLinkRepositoryInterface interface {
	Save(entity *models.SlackConversationLink) (*models.SlackConversationLink, error)
	GetBySlackIds(slackChannelId, slackThreadTs string) (*models.SlackConversationLink, error)
	GetByBrokerConversationId(brokerConversationId uuid.UUID) (*models.SlackConversationLink, error)
}
