package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
)

type CreateAgentContextRequest struct {
	AgentByCollectionTokenIdRequest
	Key string `json:"key" params:"key"`
}

func (r *CreateAgentContextRequest) Validate() error {
	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "cannot be empty")
	}

	return r.EnsureAuthenticatedWallet()
}
