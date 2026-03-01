package responses

import "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"

type GetDatabaseConnectionSetupStatusResponse struct {
	IsCompleted bool                                 `json:"isCompleted"`
	SetupStatus *enums.DatabaseConnectionSetupStatus `json:"setupStatus"`
}
