package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/google/uuid"
)

type CreateAgentContextRequest struct {
	AgentId uuid.UUID `json:"agentId"`
	Key     string    `json:"key"`
}

func (r *CreateAgentContextRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "cannot be empty")
	}

	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "cannot be empty")
	}

	return nil
}
