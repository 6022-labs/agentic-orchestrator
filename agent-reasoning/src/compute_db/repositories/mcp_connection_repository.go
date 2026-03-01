package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/repositories"
	"github.com/google/uuid"
)

type MCPConnectionRepository struct {
	client db_clients.DatabaseClientInterface
}

func NewMCPConnectionRepository(client db_clients.DatabaseClientInterface) *MCPConnectionRepository {
	return &MCPConnectionRepository{
		client: client,
	}
}

func (r *MCPConnectionRepository) Save(entity *models.MCPConnection) (*models.MCPConnection, error) {
	if err := r.client.GetDB().Save(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *MCPConnectionRepository) Delete(entity *models.MCPConnection) error {
	return r.DeleteById(entity.ID)
}

func (r *MCPConnectionRepository) DeleteById(id uuid.UUID) error {
	tx := r.client.
		GetDB().
		Where("id = ?", id).
		Delete(&models.MCPConnection{})

	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (r *MCPConnectionRepository) GetById(id uuid.UUID) (*models.MCPConnection, error) {
	var entity models.MCPConnection

	tx := r.client.
		GetDB().
		Where("id = ?", id).
		First(&entity)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *MCPConnectionRepository) GetByAgentId(agentId uuid.UUID) ([]*models.MCPConnection, error) {
	var entities []*models.MCPConnection

	tx := r.client.
		GetDB().
		Where("agent_id = ?", agentId).
		Find(&entities)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *MCPConnectionRepository) Search(params repositories.MCPConnectionSearchParams) ([]*models.MCPConnection, error) {
	var entities []*models.MCPConnection
	query := r.client.GetDB().Model(&models.MCPConnection{})

	if params.AgentId != nil {
		query = query.Where("agent_id = ?", *params.AgentId)
	}

	if params.Search != nil {
		likePattern := "%" + *params.Search + "%"
		query = query.Where("name ILIKE ? OR url ILIKE ?", likePattern, likePattern)
	}

	if params.Offset != nil {
		query = query.Offset(int(*params.Offset))
	}

	if params.Limit != nil {
		query = query.Limit(int(*params.Limit))
	}

	if params.OrderBy != nil {
		if params.OrderByDirection == nil {
			defaultDirection := "ASC"
			params.OrderByDirection = &defaultDirection
		}

		query = query.Order(*params.OrderBy + " " + *params.OrderByDirection)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *MCPConnectionRepository) ExistsByAgentIdAndUrl(agentId uuid.UUID, url string) (*bool, error) {
	var count int64

	tx := r.client.
		GetDB().
		Model(&models.MCPConnection{}).
		Where("agent_id = ? AND url = ?", agentId, url).
		Count(&count)

	if err := tx.Error; err != nil {
		return nil, err
	}

	result := count > 0
	return &result, nil
}

func (r *MCPConnectionRepository) ExistsByAgentIdAndName(agentId uuid.UUID, name string) (*bool, error) {
	var count int64

	tx := r.client.
		GetDB().
		Model(&models.MCPConnection{}).
		Where("agent_id = ? AND name = ?", agentId, name).
		Count(&count)

	if err := tx.Error; err != nil {
		return nil, err
	}

	result := count > 0
	return &result, nil
}
