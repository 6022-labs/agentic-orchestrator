package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/google/uuid"
)

type GetLinkedAgentContextOfConversationRequest struct {
	AgentId        uuid.UUID `json:"agentId" query:"agentId"`
	ConversationId string    `json:"conversationId" query:"conversationId"`
}

func (r *GetLinkedAgentContextOfConversationRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "cannot be empty")
	}

	if len(strings.TrimSpace(r.ConversationId)) == 0 {
		return errors.NewValidationError("conversationId", "cannot be empty")
	}

	return nil
}
