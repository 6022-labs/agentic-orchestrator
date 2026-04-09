package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
)

type AgentCollectionRepositoryInterface interface {
	RepositoryInterface[models.AgentCollection]
	GetAll() ([]*models.AgentCollection, error)
	GetByAddress(address string) (*models.AgentCollection, error)
	ExistsByCollectionAddress(address string) (bool, error)
}
