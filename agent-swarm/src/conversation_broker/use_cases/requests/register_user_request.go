package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
)

type RegisterUserRequest struct {
	DisplayName string `json:"displayName"`
}

func (r *RegisterUserRequest) Validate() error {
	if strings.TrimSpace(r.DisplayName) == "" {
		return errors.NewValidationError("displayName", "must not be empty")
	}
	return nil
}
