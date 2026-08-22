package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
)

type SlackAgentRepositoryInterface interface {
	Save(ctx context.Context, entity *models.SlackAgent) (*models.SlackAgent, error)
	UpdateOwnerWallet(ctx context.Context, collectionAddress string, collectionTokenId uint64, ownerWallet string) error
	UpdateIconImageCid(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
		iconImageCid *string,
	) error
	GetByCollectionAddressAndCollectionTokenId(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
	) (*models.SlackAgent, error)
}
