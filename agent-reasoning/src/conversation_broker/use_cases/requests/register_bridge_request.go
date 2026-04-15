package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
)

type RegisterBridgeRequest struct {
	Kind        enums.BridgeKind `json:"kind"`
	CallbackUrl string           `json:"callbackUrl"`
}

func (r *RegisterBridgeRequest) Validate() error {
	if !enums.GetAllBridgeKinds().Contains(r.Kind) {
		return errors.NewValidationError("kind", "must be a valid bridge kind")
	}
	if strings.TrimSpace(r.CallbackUrl) == "" {
		return errors.NewValidationError("callbackUrl", "must not be empty")
	}
	return nil
}
