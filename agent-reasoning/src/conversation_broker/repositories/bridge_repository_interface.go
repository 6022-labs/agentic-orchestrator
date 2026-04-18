package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type BridgeRepositoryInterface interface {
	Save(entity *models.Bridge) (*models.Bridge, error)
	GetById(id uuid.UUID) (*models.Bridge, error)
	ListByKind(kind enums.BridgeKind) ([]*models.Bridge, error)
}
