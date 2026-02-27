package repositories

import "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db"

type AbstractRepository[T any] struct {
	client *slack_bridge_db.DatabaseClient
}

func NewAbstractRepository[T any](
	client *slack_bridge_db.DatabaseClient,
) *AbstractRepository[T] {
	return &AbstractRepository[T]{
		client: client,
	}
}

func (r *AbstractRepository[T]) Save(entity *T) (*T, error) {
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
