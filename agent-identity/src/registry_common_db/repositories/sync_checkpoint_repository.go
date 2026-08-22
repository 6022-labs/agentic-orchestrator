package repositories

import (
	"context"

	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"gorm.io/gorm"
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

func (r *SyncCheckpointRepository) Get(ctx context.Context) (*models.SyncCheckpoint, error) {
	checkpoint := &models.SyncCheckpoint{}

	tx := r.client.GetDB().WithContext(ctx).First(checkpoint)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return checkpoint, nil
}
