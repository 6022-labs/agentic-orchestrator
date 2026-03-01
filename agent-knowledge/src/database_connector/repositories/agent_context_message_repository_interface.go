package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentContextMessageRepositoryInterface interface {
	SaveMultiple(entities []*models.AgentContextMessage) error
	CreateMultiple(entities []*models.AgentContextMessage) ([]*models.AgentContextMessage, error)
	CreateMultipleUsingTransaction(tx *gorm.DB, entities []*models.AgentContextMessage) ([]*models.AgentContextMessage, error)
	GetAll() ([]*models.AgentContextMessage, error)
	GetByAgentContextKey(agentContextKey string) ([]*models.AgentContextMessage, error)
	GetByAgentContextId(agentContextId uuid.UUID) ([]*models.AgentContextMessage, error)
}
