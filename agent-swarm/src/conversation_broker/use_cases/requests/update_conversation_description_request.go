package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/google/uuid"
)

type UpdateConversationDescriptionRequest struct {
	ConversationID uuid.UUID `json:"-"           params:"id"`
	Description    *string   `json:"description"`
}

func (r *UpdateConversationDescriptionRequest) Validate() error {
	if r.ConversationID == uuid.Nil {
		return errors.NewValidationError("conversationId", "must not be empty")
	}
	return nil
}
