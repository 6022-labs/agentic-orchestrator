package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/repositories"
)

type AgentContextMessageRepositoryFactoryInterface interface {
	Create(dbClient db_clients.AgentDatabaseClientInterface) repositories.AgentContextMessageRepositoryInterface
}
