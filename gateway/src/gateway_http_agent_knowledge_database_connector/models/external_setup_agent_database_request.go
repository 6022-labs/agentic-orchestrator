package models

import "github.com/google/uuid"

// ExternalSetupAgentDatabaseRequest represents the payload for setting up an agent database.
// Matches swagger requests.SetupAgentDatabaseRequest with optional copyFromAgentId
type ExternalSetupAgentDatabaseRequest struct {
	AgentId         uuid.UUID  `json:"agentId"`
	CopyFromAgentId *uuid.UUID `json:"copyFromAgentId,omitempty"`
}
