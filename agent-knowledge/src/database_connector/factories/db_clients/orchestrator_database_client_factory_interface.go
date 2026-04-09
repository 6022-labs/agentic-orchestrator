package db_clients

import "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"

type OrchestratorDatabaseClientFactoryInterface interface {
	Create() (db_clients.OrchestratorDatabaseClientInterface, error)
}
