package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/repositories"
)

type AgentContextConversationMappingRepositoryFactoryInterface interface {
	Create(dbClient db_clients.AgentDatabaseClientInterface) repositories.AgentContextConversationMappingRepositoryInterface
}
