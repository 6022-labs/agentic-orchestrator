package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackChannelRepositoryInterface interface {
	Save(ctx context.Context, entity *models.SlackChannel) (*models.SlackChannel, error)
	GetById(ctx context.Context, id string) (*models.SlackChannel, error)
	AddAgent(ctx context.Context, entity *models.SlackChannel, agent *models.SlackAgent) (*models.SlackChannel, error)
	RemoveAgent(
		ctx context.Context,
		entity *models.SlackChannel,
		agent *models.SlackAgent,
	) (*models.SlackChannel, error)
	AddUser(ctx context.Context, channelId string, userId string) error
	UpdateId(ctx context.Context, oldId string, newId string) error
}
