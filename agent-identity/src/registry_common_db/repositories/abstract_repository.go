package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/db_clients"
	"gorm.io/gorm"
)

type AbstractRepository[T any] struct {
	client db_clients.DatabaseClientInterface
}

func NewAbstractRepository[T any](client db_clients.DatabaseClientInterface) *AbstractRepository[T] {
	return &AbstractRepository[T]{client: client}
}

func (r *AbstractRepository[T]) Save(ctx context.Context, entity *T) (*T, error) {
	return r.SaveUsingTransaction(ctx, entity, r.client.GetDB())
}

func (r *AbstractRepository[T]) SaveUsingTransaction(ctx context.Context, entity *T, tx *gorm.DB) (*T, error) {
	if err := tx.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
