package repositories

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"
)

type SlackAgentRepository struct {
	AbstractRepository[models.SlackAgent]
}

func NewSlackAgentRepository(
	client *slack_bridge_db.DatabaseClient,
) *SlackAgentRepository {
	return &SlackAgentRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackAgent](client),
	}
}

func (r *SlackAgentRepository) GetByCollectionAddressAndCollectionTokenId(collectionAddress string, collectionTokenId uint64) (*models.SlackAgent, error) {
	var entity models.SlackAgent
	err := r.client.
		Where("collection_address = ? AND collection_token_id = ?", collectionAddress, collectionTokenId).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *SlackAgentRepository) GetByEnsDomain(ensDomain string) (*models.SlackAgent, error) {
	var entity models.SlackAgent
	err := r.client.
		Where("ens_domain = ?", ensDomain).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *SlackAgentRepository) GetNextAgentToProcess(processingInterval uint64) (*models.SlackAgent, error) {
	var entity models.SlackAgent
	err := r.client.
		Where("is_processing = ?", false).
		Where(
			fmt.Sprintf(
				"last_processed_at IS NULL OR last_processed_at < NOW() - INTERVAL '%d seconds'",
				processingInterval,
			),
		).
		Order("last_processed_at ASC NULLS FIRST").
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *SlackAgentRepository) UpdateOwnerWallet(collectionAddress string, collectionTokenId uint64, ownerWallet string) error {
	result := r.client.Model(&models.SlackAgent{}).
		Where("collection_address = ? AND collection_token_id = ?", collectionAddress, collectionTokenId).
		Update("owner_wallet", ownerWallet)

	if result.RowsAffected == 0 {
		return fmt.Errorf("error while updating owner wallet, no rows updated for collection_address %s and collection_token_id %d", collectionAddress, collectionTokenId)
	}

	return result.Error
}

func (r *SlackAgentRepository) UpdateIconImageCid(collectionAddress string, collectionTokenId uint64, iconImageCid *string) error {
	result := r.client.Model(&models.SlackAgent{}).
		Where("collection_address = ? AND collection_token_id = ?", collectionAddress, collectionTokenId).
		Update("icon_image_cid", iconImageCid)

	if result.RowsAffected == 0 {
		return fmt.Errorf("error while updating icon image cid, no rows updated for collection_address %s and collection_token_id %d", collectionAddress, collectionTokenId)
	}

	return result.Error
}
