package repositories

import (
	"context"

	"errors"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *ParticipantRepository) ClaimNextAgentToProcess(
	ctx context.Context,
	processingIntervalSeconds uint64,
	leaseDuration time.Duration,
) (*models.Participant, error) {
	db := r.client.GetDB().WithContext(ctx)
	next := db.Model(&models.Participant{}).
		Select("id").
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("collection_address IS NOT NULL AND collection_token_id IS NOT NULL").
		Where("lease_until IS NULL OR lease_until < NOW()").
		Where(
			"last_processed_at IS NULL OR last_processed_at < NOW() - make_interval(secs => ?)",
			processingIntervalSeconds,
		).
		Order("last_processed_at ASC NULLS FIRST").
		Limit(1)

	var claimed []*models.Participant
	err := db.Model(&claimed).
		Clauses(clause.Returning{}).
		Where("id = (?)", next).
		Updates(map[string]any{"lease_id": uuid.New(), "lease_until": r.leaseUntil(leaseDuration)}).
		Error
	if err != nil {
		return nil, err
	}
	if len(claimed) == 0 {
		return nil, nil
	}
	return claimed[0], nil
}

func (r *ParticipantRepository) RenewLease(
	ctx context.Context,
	participant *models.Participant,
	leaseDuration time.Duration,
) (bool, error) {
	result := r.leased(ctx, participant).Update("lease_until", r.leaseUntil(leaseDuration))
	return result.RowsAffected == 1, result.Error
}

func (r *ParticipantRepository) ReleaseLease(ctx context.Context, participant *models.Participant) error {
	return r.leased(ctx, participant).
		Updates(map[string]any{"lease_id": nil, "lease_until": nil, "last_processed_at": gorm.Expr("NOW()")}).
		Error
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

func (r *ParticipantRepository) leased(ctx context.Context, participant *models.Participant) *gorm.DB {
	return r.client.GetDB().WithContext(ctx).
		Model(&models.Participant{}).
		Where("id = ? AND lease_id = ?", participant.ID, participant.LeaseID)
}

// leaseUntil uses the database clock, which every replica shares.
func (r *ParticipantRepository) leaseUntil(leaseDuration time.Duration) clause.Expr {
	return gorm.Expr("NOW() + make_interval(secs => ?)", leaseDuration.Seconds())
}
