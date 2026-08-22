package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
)

type AgentCollectionRepositoryInterface interface {
	RepositoryInterface[models.AgentCollection]
	GetAll(ctx context.Context) ([]*models.AgentCollection, error)
	GetByAddress(ctx context.Context, address string) (*models.AgentCollection, error)
	ExistsByCollectionAddress(ctx context.Context, address string) (bool, error)
}
