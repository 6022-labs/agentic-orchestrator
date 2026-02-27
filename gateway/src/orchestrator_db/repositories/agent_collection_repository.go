package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/db_clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
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

func (r *AgentCollectionRepository) GetAll() ([]*models.AgentCollection, error) {
	collections := []*models.AgentCollection{}

	tx := r.client.GetDB().Find(&collections)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return collections, nil
}

func (r *AgentCollectionRepository) GetByAddress(address string) (*models.AgentCollection, error) {
	collection := &models.AgentCollection{}

	tx := r.client.GetDB().Where("address = ?", address).First(collection)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return collection, nil
}

func (r *AgentCollectionRepository) ExistsByCollectionAddress(address string) (bool, error) {
	var count int64

	tx := r.client.GetDB().Model(&models.AgentCollection{}).Where("address = ?", address).Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}

	return count > 0, nil
}
