package repositories

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ParticipantRepository struct {
	AbstractRepository[models.Participant]
}

func NewParticipantRepository(client db_clients.DatabaseClientInterface) *ParticipantRepository {
	return &ParticipantRepository{
		AbstractRepository: *NewAbstractRepository[models.Participant](client),
	}
}

func (r *ParticipantRepository) GetById(id uuid.UUID) (*models.Participant, error) {
	var entity models.Participant
	err := r.client.GetDB().
		Where("id = ?", id).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ParticipantRepository) GetByCollectionAddressAndCollectionTokenId(
	collectionAddress string,
	collectionTokenId uint64,
) (*models.Participant, error) {
	var entity models.Participant
	err := r.client.GetDB().
		Where("collection_address = ? AND collection_token_id = ?", collectionAddress, collectionTokenId).
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetNextMintedParticipantToProcess returns the oldest minted participant
// (role IS NOT NULL) that is not currently processing and whose last run is
// older than processingIntervalSeconds. Non-minted humans are never returned.
func (r *ParticipantRepository) GetNextMintedParticipantToProcess(
	processingIntervalSeconds uint64,
) (*models.Participant, error) {
	var entity models.Participant
	err := r.client.GetDB().
		Where("role IS NOT NULL").
		Where("is_processing = ?", false).
		Where(fmt.Sprintf(
			"last_processed_at IS NULL OR last_processed_at < NOW() - INTERVAL '%d seconds'",
			processingIntervalSeconds,
		)).
		Order("last_processed_at ASC NULLS FIRST").
		First(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
