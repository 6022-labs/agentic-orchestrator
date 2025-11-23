package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentContextSearchParams struct {
	Search        *string
	SortBy        *string
	SortDirection *string
}

type AgentContextRepositoryInterface interface {
	Save(entity *models.AgentContext) (*models.AgentContext, error)
	CreateMultiple(entities []*models.AgentContext) ([]*models.AgentContext, error)
	CreateMultipleUsingTransaction(tx *gorm.DB, entities []*models.AgentContext) ([]*models.AgentContext, error)
	GetById(id uuid.UUID) (*models.AgentContext, error)
	GetByKey(key string) (*models.AgentContext, error)
	GetAll() ([]*models.AgentContext, error)
	Search(params AgentContextSearchParams) ([]*models.AgentContext, error)
}
