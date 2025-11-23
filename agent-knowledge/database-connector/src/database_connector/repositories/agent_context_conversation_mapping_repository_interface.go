package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"gorm.io/gorm"
)

type AgentContextConversationMappingRepositoryInterface interface {
	Save(entity *models.AgentContextConversationMapping) (*models.AgentContextConversationMapping, error)
	CreateMultiple(entities []*models.AgentContextConversationMapping) ([]*models.AgentContextConversationMapping, error)
	CreateMultipleUsingTransaction(tx *gorm.DB, entities []*models.AgentContextConversationMapping) ([]*models.AgentContextConversationMapping, error)
	GetAll() ([]*models.AgentContextConversationMapping, error)
	GetByConversationId(conversationId string) (*models.AgentContextConversationMapping, error)
}
