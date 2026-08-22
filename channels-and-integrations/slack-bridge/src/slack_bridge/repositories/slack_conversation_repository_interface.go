package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/google/uuid"
)

type SlackConversationRepositoryInterface interface {
	Save(ctx context.Context, entity *models.SlackConversation) (*models.SlackConversation, error)
	GetBySlackIds(ctx context.Context, slackChannelId, slackThreadTs string) (*models.SlackConversation, error)
	GetByBrokerConversationId(ctx context.Context, brokerConversationId uuid.UUID) (*models.SlackConversation, error)
	ListByChannelId(ctx context.Context, slackChannelId string) ([]*models.SlackConversation, error)
}
