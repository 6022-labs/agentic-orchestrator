package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	db_agent_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/repositories"
)

type AgentContextMessageRepositoryFactory struct {
}

func NewAgentContextMessageRepositoryFactory() *AgentContextMessageRepositoryFactory {
	return &AgentContextMessageRepositoryFactory{}
}

func (f *AgentContextMessageRepositoryFactory) Create(dbClient db_clients.AgentDatabaseClientInterface) repositories.AgentContextMessageRepositoryInterface {
	return db_agent_repositories.NewAgentContextMessageRepository(dbClient)
}
