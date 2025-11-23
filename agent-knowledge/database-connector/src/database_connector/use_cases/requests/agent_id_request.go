package requests

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/google/uuid"
)

type AgentIdRequest struct {
	AgentId uuid.UUID `json:"agentId" query:"agentId"`
}

func (r *AgentIdRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "Agent ID cannot be empty")
	}

	return nil
}
