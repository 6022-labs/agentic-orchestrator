package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/google/uuid"
)

type GetAgentContextByKeyRequest struct {
	AgentId uuid.UUID `params:"agentId"`
	Key     string    `params:"key"`
}

func (r *GetAgentContextByKeyRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "agentId is required and must be a valid UUID")
	}

	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "key is required")
	}

	return nil
}
