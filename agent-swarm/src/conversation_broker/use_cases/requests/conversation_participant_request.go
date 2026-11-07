package requests

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/google/uuid"
)

// ConversationParticipantRequest is shared by Attach and Detach use cases.
type ConversationParticipantRequest struct {
	ConversationID uuid.UUID `json:"-"             uri:"id"`
	ParticipantID  uuid.UUID `json:"participantId" uri:"participantId"`
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
