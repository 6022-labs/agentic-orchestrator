package requests

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/errors"
	"github.com/google/uuid"
)

const conversationContextDiscoveryKeyPrefix = "conversation-context-discovery"

type CreateAgentContextRequest struct {
	AgentId    uuid.UUID `json:"agentId"`
	Key        string    `json:"key"`
	IsInternal bool      `json:"isInternal"`
}

func (r *CreateAgentContextRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "cannot be empty")
	}

	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "cannot be empty")
	}

	if !r.IsInternal && strings.HasPrefix(r.Key, conversationContextDiscoveryKeyPrefix) {
		return errors.NewValidationError("key", fmt.Sprintf("cannot start with %s", conversationContextDiscoveryKeyPrefix))
	}

	return nil
}
