package requests

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/compute"
)

type CreateAgentContextRequest struct {
	AgentByCollectionTokenIdRequest
	Key string `json:"key" params:"key"`
}

func (r *CreateAgentContextRequest) Validate() error {
	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "cannot be empty")
	}

	if strings.HasPrefix(r.Key, compute.ConversationContextDiscoveryKeyPrefix) {
		return errors.NewValidationError("key", fmt.Sprintf("cannot start with %s", compute.ConversationContextDiscoveryKeyPrefix))
	}

	return r.EnsureAuthenticatedWallet()
}
