package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/repositories"
	db_agent_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector_db_agent/repositories"
)

type AgentContextRepositoryFactory struct {
}

func NewAgentContextRepositoryFactory() *AgentContextRepositoryFactory {
	return &AgentContextRepositoryFactory{}
}

func (f *AgentContextRepositoryFactory) Create(dbClient db_clients.AgentDatabaseClientInterface) repositories.AgentContextRepositoryInterface {
	return db_agent_repositories.NewAgentContextRepository(dbClient)
}
