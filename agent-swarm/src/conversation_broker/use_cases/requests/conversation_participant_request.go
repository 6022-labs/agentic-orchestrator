package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/google/uuid"
)

// ConversationParticipantRequest is shared by Attach and Detach use cases.
type ConversationParticipantRequest struct {
	ConversationID uuid.UUID `json:"-"             params:"id"`
	ParticipantID  uuid.UUID `json:"participantId" params:"participantId"`
}

func (r *ConversationParticipantRequest) Validate() error {
	if r.ConversationID == uuid.Nil {
		return errors.NewValidationError("conversationId", "must not be empty")
	}
	if r.ParticipantID == uuid.Nil {
		return errors.NewValidationError("participantId", "must not be empty")
	}
	return nil
}
