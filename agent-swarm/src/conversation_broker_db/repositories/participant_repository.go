package repositories

import (
	"context"

	"errors"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	AbstractRepository[models.Participant]
}

func NewParticipantRepository(client db_clients.DatabaseClientInterface) *ParticipantRepository {
	return &ParticipantRepository{
		AbstractRepository: *NewAbstractRepository[models.Participant](client),
	}
}

func (r *ParticipantRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Participant, error) {
	var entity models.Participant
	err := r.client.GetDB().WithContext(ctx).
		Where("id = ?", id).
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ParticipantRepository) GetByCollectionAddressAndCollectionTokenId(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
) (*models.Participant, error) {
	var entity models.Participant
	err := r.client.GetDB().WithContext(ctx).
		Where("collection_address = ? AND collection_token_id = ?", collectionAddress, collectionTokenId).
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetNextAgentToProcess returns the oldest agent participant (role IS NOT
// NULL) that is not currently processing and whose last run is older than
// processingIntervalSeconds. User participants (display-name only) are never
// returned.
func (r *ParticipantRepository) GetNextAgentToProcess(
	ctx context.Context,
	processingIntervalSeconds uint64,
) (*models.Participant, error) {
	var entity models.Participant
	err := r.client.GetDB().WithContext(ctx).
		Where("collection_address IS NOT NULL AND collection_token_id IS NOT NULL").
		Where("is_processing = ?", false).
		Where(fmt.Sprintf(
			"last_processed_at IS NULL OR last_processed_at < NOW() - INTERVAL '%d seconds'",
			processingIntervalSeconds,
		)).
		Order("last_processed_at ASC NULLS FIRST").
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *ParticipantRepository) ListByConversationId(
	ctx context.Context,
	conversationId uuid.UUID,
) ([]*models.Participant, error) {
	var entities []*models.Participant
	err := r.client.GetDB().WithContext(ctx).
		Joins("INNER JOIN conversation_participants ON conversation_participants.participant_id = participants.id").
		Where("conversation_participants.conversation_id = ?", conversationId).
		Order("conversation_participants.joined_at ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
