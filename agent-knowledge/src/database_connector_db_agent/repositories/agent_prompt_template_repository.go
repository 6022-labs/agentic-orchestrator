package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
)

type AgentPromptTemplateRepository struct {
	*AbstractRepository[models.AgentPromptTemplate]
}

func NewAgentPromptTemplateRepository(client db_clients.AgentDatabaseClientInterface) *AgentPromptTemplateRepository {
	return &AgentPromptTemplateRepository{
		AbstractRepository: NewAbstractRepository[models.AgentPromptTemplate](client),
	}
}

func (r *AgentPromptTemplateRepository) GetByKey(key enums.AgentPromptTemplateKey) (*models.AgentPromptTemplate, error) {
	var agentPromptTemplate models.AgentPromptTemplate

	quotedColumn := quoteColumnName(r.client.GetDriver(), "key")
	if err := r.client.GetDB().
		Where(quotedColumn+" = ?", key).
		First(&agentPromptTemplate).Error; err != nil {
		return nil, err
	}

	return &agentPromptTemplate, nil
}
