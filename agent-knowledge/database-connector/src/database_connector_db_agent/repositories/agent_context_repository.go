package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/google/uuid"
)

type AgentContextRepository struct {
	*AbstractRepository[models.AgentContext]
}

func NewAgentContextRepository(client db_clients.AgentDatabaseClientInterface) *AgentContextRepository {
	return &AgentContextRepository{
		AbstractRepository: NewAbstractRepository[models.AgentContext](client),
	}
}

func (r *AgentContextRepository) GetById(id uuid.UUID) (*models.AgentContext, error) {
	var entity models.AgentContext

	quotedColumn := quoteColumnName(r.client.GetDriver(), "id")
	if err := r.client.GetDB().Where(quotedColumn+" = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *AgentContextRepository) GetByKey(key string) (*models.AgentContext, error) {
	var entity models.AgentContext

	quotedColumn := quoteColumnName(r.client.GetDriver(), "key")
	if err := r.client.GetDB().Where(quotedColumn+" = ?", key).First(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *AgentContextRepository) GetAll() ([]*models.AgentContext, error) {
	var agentContexts []*models.AgentContext

	if err := r.client.GetDB().Find(&agentContexts).Error; err != nil {
		return nil, err
	}

	return agentContexts, nil
}

func (r *AgentContextRepository) Search(params repositories.AgentContextSearchParams) ([]*models.AgentContext, error) {
	var agentContexts []*models.AgentContext

	query := r.client.GetDB().Model(&models.AgentContext{})

	if params.Search != nil {
		quotedColumn := quoteColumnName(r.client.GetDriver(), "key")
		query = query.Where(ilikeContainsPattern(r.client.GetDriver(), quotedColumn), *params.Search)
	}

	if params.SortBy != nil {
		sortByQuoted := quoteColumnName(r.client.GetDriver(), *params.SortBy)
		query = query.Order(sortByQuoted + " " + *params.SortDirection)
	}

	if err := query.Find(&agentContexts).Error; err != nil {
		return nil, err
	}

	return agentContexts, nil
}
