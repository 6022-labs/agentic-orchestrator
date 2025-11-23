package responses

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"

type GetDatabaseConnectionSetupStatusResponse struct {
	IsCompleted bool                                 `json:"isCompleted"`
	SetupStatus *enums.DatabaseConnectionSetupStatus `json:"setupStatus"`
}
