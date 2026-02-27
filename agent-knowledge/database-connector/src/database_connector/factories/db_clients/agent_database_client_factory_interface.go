package db_clients

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/google/uuid"
)

type AgentDatabaseClientFactoryInterface interface {
	Create(agentId uuid.UUID) (db_clients.AgentDatabaseClientInterface, error)
	CreateFromDatabaseConnection(databaseConnection *models.DatabaseConnection) (db_clients.AgentDatabaseClientInterface, error)
}
