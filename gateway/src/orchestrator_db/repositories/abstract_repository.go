package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/db_clients"
	"gorm.io/gorm"
)

type AbstractRepository[T any] struct {
	client db_clients.DatabaseClientInterface
}

func NewAbstractRepository[T any](client db_clients.DatabaseClientInterface) *AbstractRepository[T] {
	return &AbstractRepository[T]{client: client}
}

func (r *AbstractRepository[T]) Save(entity *T) (*T, error) {
	return r.SaveUsingTransaction(entity, r.client.GetDB())
}

func (r *AbstractRepository[T]) SaveUsingTransaction(entity *T, tx *gorm.DB) (*T, error) {
	if err := tx.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
