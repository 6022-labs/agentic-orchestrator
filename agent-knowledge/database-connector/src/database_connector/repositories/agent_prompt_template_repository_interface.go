package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"gorm.io/gorm"
)

type AgentPromptTemplateRepositoryInterface interface {
	Save(agentPromptTemplate *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error)
	CreateMultiple(entities []*models.AgentPromptTemplate) ([]*models.AgentPromptTemplate, error)
	CreateMultipleUsingTransaction(tx *gorm.DB, entities []*models.AgentPromptTemplate) ([]*models.AgentPromptTemplate, error)
	GetByKey(key enums.AgentPromptTemplateKey) (*models.AgentPromptTemplate, error)
	GetAll() ([]*models.AgentPromptTemplate, error)
}
