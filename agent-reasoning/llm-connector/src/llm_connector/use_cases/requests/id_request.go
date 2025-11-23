package requests

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/errors"
	"github.com/google/uuid"
)

type IdRequest struct {
	ID uuid.UUID `json:"id"`
}

func (r *IdRequest) Validate() error {
	if r.ID == uuid.Nil {
		return errors.NewValidationError("id", "ID cannot be empty")
	}

	return nil
}
