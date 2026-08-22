package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
)

type AbstractRepository[T any] struct {
	client db_clients.DatabaseClientInterface
}

func NewAbstractRepository[T any](client db_clients.DatabaseClientInterface) *AbstractRepository[T] {
	return &AbstractRepository[T]{client: client}
}

func (r *AbstractRepository[T]) Save(ctx context.Context, entity *T) (*T, error) {
	if err := r.client.GetDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
