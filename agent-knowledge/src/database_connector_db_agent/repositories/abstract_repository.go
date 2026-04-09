package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"gorm.io/gorm"
)

type AbstractRepository[TModel any] struct {
	client db_clients.AgentDatabaseClientInterface
}

func NewAbstractRepository[TModel any](client db_clients.AgentDatabaseClientInterface) *AbstractRepository[TModel] {
	return &AbstractRepository[TModel]{
		client: client,
	}
}

func (r *AbstractRepository[TModel]) Save(entity *TModel) (*TModel, error) {
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *AbstractRepository[TModel]) GetAll() ([]*TModel, error) {
	var entities []*TModel
	if err := r.client.GetDB().Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *AbstractRepository[TModel]) CreateMultiple(entities []*TModel) ([]*TModel, error) {
	return r.CreateMultipleUsingTransaction(r.client.GetDB(), entities)
}

func (r *AbstractRepository[TModel]) CreateMultipleUsingTransaction(tx *gorm.DB, entities []*TModel) ([]*TModel, error) {
	if err := tx.CreateInBatches(entities, 100).Error; err != nil {
		return nil, err
	}
	return entities, nil
}
