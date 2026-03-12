package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
)

type SyncCheckpointRepositoryInterface interface {
	repositories.RepositoryInterface[models.SyncCheckpoint]
	Get() (*models.SyncCheckpoint, error)
}
