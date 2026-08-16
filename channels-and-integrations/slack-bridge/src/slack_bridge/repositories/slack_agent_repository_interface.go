package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackAgentRepositoryInterface interface {
	Save(entity *models.SlackAgent) (*models.SlackAgent, error)
	UpdateOwnerWallet(collectionAddress string, collectionTokenId uint64, ownerWallet string) error
	UpdateIconImageCid(collectionAddress string, collectionTokenId uint64, iconImageCid *string) error
	GetByCollectionAddressAndCollectionTokenId(
		collectionAddress string,
		collectionTokenId uint64,
	) (*models.SlackAgent, error)
}
