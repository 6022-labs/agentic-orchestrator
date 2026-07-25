package repositories

import (
	"errors"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BridgeRepository struct {
	AbstractRepository[models.Bridge]
}

func NewBridgeRepository(client db_clients.DatabaseClientInterface) *BridgeRepository {
	return &BridgeRepository{
		AbstractRepository: *NewAbstractRepository[models.Bridge](client),
	}
}

func (r *BridgeRepository) GetById(id uuid.UUID) (*models.Bridge, error) {
	var entity models.Bridge
	err := r.client.GetDB().
		Where("id = ?", id).
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BridgeRepository) ListByKind(kind enums.BridgeKind) ([]*models.Bridge, error) {
	var entities []*models.Bridge
	err := r.client.GetDB().
		Where("kind = ?", kind).
		Order("created_at ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
