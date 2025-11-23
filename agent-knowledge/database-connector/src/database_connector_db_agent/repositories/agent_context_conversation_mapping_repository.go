package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
)

type AgentContextConversationMappingRepository struct {
	*AbstractRepository[models.AgentContextConversationMapping]
}

func NewAgentContextConversationMappingRepository(client db_clients.AgentDatabaseClientInterface) *AgentContextConversationMappingRepository {
	return &AgentContextConversationMappingRepository{
		AbstractRepository: NewAbstractRepository[models.AgentContextConversationMapping](client),
	}
}

func (r *AgentContextConversationMappingRepository) GetByConversationId(conversationId string) (*models.AgentContextConversationMapping, error) {
	var entity models.AgentContextConversationMapping
	if err := r.client.GetDB().Where("conversation_id = ?", conversationId).First(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}
