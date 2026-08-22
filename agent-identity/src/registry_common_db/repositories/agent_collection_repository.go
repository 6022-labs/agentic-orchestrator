package repositories

import (
	"context"

	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"gorm.io/gorm"
)

type AgentCollectionRepository struct {
	*AbstractRepository[models.AgentCollection]
	client db_clients.DatabaseClientInterface
}

func NewAgentCollectionRepository(client db_clients.DatabaseClientInterface) *AgentCollectionRepository {
	abstractRepository := NewAbstractRepository[models.AgentCollection](client)

	return &AgentCollectionRepository{
		client:             client,
		AbstractRepository: abstractRepository,
	}
}

func (r *AgentCollectionRepository) GetAll(ctx context.Context) ([]*models.AgentCollection, error) {
	collections := []*models.AgentCollection{}

	tx := r.client.GetDB().WithContext(ctx).Find(&collections)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return collections, nil
}

func (r *AgentCollectionRepository) GetByAddress(ctx context.Context, address string) (*models.AgentCollection, error) {
	collection := &models.AgentCollection{}

	tx := r.client.GetDB().WithContext(ctx).Where("address = ?", address).First(collection)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return collection, nil
}

func (r *AgentCollectionRepository) ExistsByCollectionAddress(ctx context.Context, address string) (bool, error) {
	var count int64

	tx := r.client.GetDB().WithContext(ctx).Model(&models.AgentCollection{}).Where("address = ?", address).Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}

	return count > 0, nil
}
