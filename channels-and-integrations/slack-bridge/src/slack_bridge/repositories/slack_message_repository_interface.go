package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackMessageRepositoryInterface interface {
	Save(ctx context.Context, entity *models.SlackMessage) (*models.SlackMessage, error)
	GetBySlackIds(ctx context.Context, slackChannelId, slackMessageTs string) (*models.SlackMessage, error)
}
