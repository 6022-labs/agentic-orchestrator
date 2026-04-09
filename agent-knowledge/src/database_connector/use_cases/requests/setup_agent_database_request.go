package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/errors"
	"github.com/google/uuid"
)

type SetupAgentDatabaseRequest struct {
	AgentId         uuid.UUID  `json:"agentId"`
	CopyFromAgentId *uuid.UUID `json:"copyFromAgentId,omitempty"`
}

func (r *SetupAgentDatabaseRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "is required")
	}

	if r.CopyFromAgentId != nil && *r.CopyFromAgentId == uuid.Nil {
		return errors.NewValidationError("copyFromAgentId", "if provided, must be a valid UUID")
	}

	if r.CopyFromAgentId != nil && *r.CopyFromAgentId == r.AgentId {
		return errors.NewValidationError("copyFromAgentId", "cannot be the same as agentId")
	}

	return nil
}
