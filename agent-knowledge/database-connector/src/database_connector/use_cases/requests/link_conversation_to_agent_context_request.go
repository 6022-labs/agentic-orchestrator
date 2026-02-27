package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/google/uuid"
)

type LinkConversationToAgentContextRequest struct {
	Key            string    `json:"key"`
	AgentId        uuid.UUID `json:"agentId"`
	ConversationId string    `json:"conversationId"`
}

func (r *LinkConversationToAgentContextRequest) Validate() error {
	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "cannot be empty")
	}

	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "cannot be empty")
	}

	if len(strings.TrimSpace(r.ConversationId)) == 0 {
		return errors.NewValidationError("conversationId", "cannot be empty")
	}

	return nil
}
