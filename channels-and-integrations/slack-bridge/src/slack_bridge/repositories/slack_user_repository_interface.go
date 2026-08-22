package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackUserRepositoryInterface interface {
	Save(ctx context.Context, entity *models.SlackUser) (*models.SlackUser, error)
	GetById(ctx context.Context, id string) (*models.SlackUser, error)
	ListByIds(ctx context.Context, ids []string) ([]*models.SlackUser, error)
	GetByPairedHumanAgent(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
	) (*models.SlackUser, error)
}
