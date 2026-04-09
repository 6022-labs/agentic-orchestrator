package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
	db_agent_repositories "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/repositories"
)

type AgentContextConversationMappingRepositoryFactory struct {
}

func NewAgentContextConversationMappingRepositoryFactory() *AgentContextConversationMappingRepositoryFactory {
	return &AgentContextConversationMappingRepositoryFactory{}
}

func (f *AgentContextConversationMappingRepositoryFactory) Create(dbClient db_clients.AgentDatabaseClientInterface) repositories.AgentContextConversationMappingRepositoryInterface {
	return db_agent_repositories.NewAgentContextConversationMappingRepository(dbClient)
}
