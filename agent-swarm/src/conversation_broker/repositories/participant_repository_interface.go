package repositories

import (
	"context"

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
	// ClaimNextAgentToProcess marks the agent least recently processed as processing and returns it, nil when none
	// is due; two replicas never claim the same one.
	ClaimNextAgentToProcess(ctx context.Context, processingIntervalSeconds uint64) (*models.Participant, error)
	ListByConversationId(ctx context.Context, conversationId uuid.UUID) ([]*models.Participant, error)
}
