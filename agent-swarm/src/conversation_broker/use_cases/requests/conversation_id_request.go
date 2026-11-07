package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/google/uuid"
)

// ConversationIdRequest is shared by ForceConclude and GetConversation.
type ConversationIdRequest struct {
	ConversationID uuid.UUID `uri:"id"`
}

func (r *ConversationIdRequest) Validate() error {
	if r.ConversationID == uuid.Nil {
		return errors.NewValidationError("conversationId", "must not be empty")
	}
	return nil
}
