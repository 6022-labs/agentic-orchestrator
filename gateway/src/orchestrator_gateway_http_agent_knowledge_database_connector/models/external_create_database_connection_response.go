package models

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"

type ExternalCreateDatabaseConnectionResponse struct {
	DatabaseName     string `json:"databaseName"`
	DatabaseUsername string `json:"databaseUsername"`
	DatabasePassword string `json:"databasePassword"`
}

func (m *ExternalCreateDatabaseConnectionResponse) ToCreateDatabaseConnectionResult() *interfaces.CreateDatabaseConnectionResult {
	return &interfaces.CreateDatabaseConnectionResult{
		DatabaseName:     m.DatabaseName,
		DatabaseUsername: m.DatabaseUsername,
		DatabasePassword: m.DatabasePassword,
	}
}
