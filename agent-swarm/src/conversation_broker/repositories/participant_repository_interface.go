package repositories

import (
	"context"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ParticipantRepositoryInterface interface {
	Save(ctx context.Context, entity *models.Participant) (*models.Participant, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Participant, error)
	GetByCollectionAddressAndCollectionTokenId(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
	) (*models.Participant, error)
	// ClaimNextAgentToProcess leases the agent least recently processed and returns it, nil when none is due; two
	// replicas never claim the same one, and an expired lease is free again.
	ClaimNextAgentToProcess(
		ctx context.Context,
		processingIntervalSeconds uint64,
		leaseDuration time.Duration,
	) (*models.Participant, error)
	// RenewLease extends the lease the participant was claimed with; false when it was lost.
	RenewLease(ctx context.Context, participant *models.Participant, leaseDuration time.Duration) (bool, error)
	// ReleaseLease frees the participant and records it as processed, unless its lease was lost.
	ReleaseLease(ctx context.Context, participant *models.Participant) error
	ListByConversationId(ctx context.Context, conversationId uuid.UUID) ([]*models.Participant, error)
}
