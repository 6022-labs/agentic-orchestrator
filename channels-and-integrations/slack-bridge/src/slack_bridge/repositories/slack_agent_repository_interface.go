package repositories

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"

type SlackAgentRepositoryInterface interface {
	GetNextAgentToProcess() (*models.SlackAgent, error)
	Save(entity *models.SlackAgent) (*models.SlackAgent, error)
	UpdateOwnerWallet(collectionAddress string, collectionTokenId uint64, ownerWallet string) error
	UpdateIconImageCid(collectionAddress string, collectionTokenId uint64, iconImageCid *string) error
	GetByEnsDomain(ensDomain string) (*models.SlackAgent, error)
	GetByCollectionAddressAndCollectionTokenId(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error)
}
