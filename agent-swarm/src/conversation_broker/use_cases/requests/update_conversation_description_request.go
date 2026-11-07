package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/google/uuid"
)

type UpdateConversationDescriptionRequest struct {
	ConversationID uuid.UUID `json:"-"           uri:"id"`
	Description    *string   `json:"description"`
}

func (r *UpdateConversationDescriptionRequest) Validate() error {
	if r.ConversationID == uuid.Nil {
		return errors.NewValidationError("conversationId", "must not be empty")
	}
	return nil
}
