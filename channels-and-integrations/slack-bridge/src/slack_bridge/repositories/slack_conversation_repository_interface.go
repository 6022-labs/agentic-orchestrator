package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/google/uuid"
)

type SlackConversationRepositoryInterface interface {
	Save(entity *models.SlackConversation) (*models.SlackConversation, error)
	GetBySlackIds(slackChannelId, slackThreadTs string) (*models.SlackConversation, error)
	GetByBrokerConversationId(brokerConversationId uuid.UUID) (*models.SlackConversation, error)
	ListByChannelId(slackChannelId string) ([]*models.SlackConversation, error)
}
