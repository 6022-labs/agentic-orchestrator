package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
)

// ResolveAgentUrlRequest identifies the agent by its ENS domain.
type ResolveAgentUrlRequest struct {
	EnsDomain string `params:"ensDomain"`
}

func (r *ResolveAgentUrlRequest) Validate() error {
	if strings.TrimSpace(r.EnsDomain) == "" {
		return errors.NewValidationError("ensDomain", "must not be empty")
	}
	return nil
}
