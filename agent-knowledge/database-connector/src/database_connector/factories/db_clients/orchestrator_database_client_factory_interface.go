package db_clients

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/db_clients"

type OrchestratorDatabaseClientFactoryInterface interface {
	Create() (db_clients.OrchestratorDatabaseClientInterface, error)
}
