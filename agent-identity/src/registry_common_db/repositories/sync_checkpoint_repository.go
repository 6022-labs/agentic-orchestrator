package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
)

type SyncCheckpointRepository struct {
	*AbstractRepository[models.SyncCheckpoint]
	client db_clients.DatabaseClientInterface
}

func NewSyncCheckpointRepository(client db_clients.DatabaseClientInterface) *SyncCheckpointRepository {
	abstractRepository := NewAbstractRepository[models.SyncCheckpoint](client)

	return &SyncCheckpointRepository{
		client:             client,
		AbstractRepository: abstractRepository,
	}
}

func (r *SyncCheckpointRepository) Get() (*models.SyncCheckpoint, error) {
	checkpoint := &models.SyncCheckpoint{}

	tx := r.client.GetDB().First(checkpoint)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return checkpoint, nil
}
