package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/google/uuid"
)

// ConversationIdRequest is shared by ForceConclude and GetConversation.
type ConversationIdRequest struct {
	ConversationID uuid.UUID `params:"id"`
}

func (r *ConversationIdRequest) Validate() error {
	if r.ConversationID == uuid.Nil {
		return errors.NewValidationError("conversationId", "must not be empty")
	}
	return nil
}
