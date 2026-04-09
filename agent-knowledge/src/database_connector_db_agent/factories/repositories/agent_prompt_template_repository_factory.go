package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	db_agent_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/repositories"
)

type AgentPromptTemplateRepositoryFactory struct {
}

func NewAgentPromptTemplateRepositoryFactory() *AgentPromptTemplateRepositoryFactory {
	return &AgentPromptTemplateRepositoryFactory{}
}

func (f *AgentPromptTemplateRepositoryFactory) Create(dbClient db_clients.AgentDatabaseClientInterface) repositories.AgentPromptTemplateRepositoryInterface {
	return db_agent_repositories.NewAgentPromptTemplateRepository(dbClient)
}
