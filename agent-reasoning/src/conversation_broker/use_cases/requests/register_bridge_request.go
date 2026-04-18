package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

// RegisterBridgeRequest registers (or upserts) a bridge with the broker.
// When ID is provided, the broker is expected to reuse the existing row if
// already known, or create it with that id otherwise. When ID is nil, the
// broker generates a fresh id.
type RegisterBridgeRequest struct {
	ID      *uuid.UUID       `json:"id,omitempty"`
	Kind    enums.BridgeKind `json:"kind"`
	BaseUrl string           `json:"baseUrl"`
}

func (r *RegisterBridgeRequest) Validate() error {
	if !enums.GetAllBridgeKinds().Contains(r.Kind) {
		return errors.NewValidationError("kind", "must be a valid bridge kind")
	}
	if strings.TrimSpace(r.BaseUrl) == "" {
		return errors.NewValidationError("baseUrl", "must not be empty")
	}
	if r.ID != nil && *r.ID == uuid.Nil {
		return errors.NewValidationError("id", "must not be the zero UUID")
	}
	return nil
}
