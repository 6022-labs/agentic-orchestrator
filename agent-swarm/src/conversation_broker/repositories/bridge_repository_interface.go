package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type BridgeRepositoryInterface interface {
	Save(ctx context.Context, entity *models.Bridge) (*models.Bridge, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Bridge, error)
	ListByKind(ctx context.Context, kind enums.BridgeKind) ([]*models.Bridge, error)
}
