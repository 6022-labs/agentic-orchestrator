package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/google/uuid"
)

type IngestMessageRequest struct {
	ConversationID uuid.UUID `json:"-" params:"id"`
	ParticipantID  uuid.UUID `json:"participantId"`
	Content        string    `json:"content"`
}

func (r *IngestMessageRequest) Validate() error {
	if r.ConversationID == uuid.Nil {
		return errors.NewValidationError("conversationId", "must not be empty")
	}
	if r.ParticipantID == uuid.Nil {
		return errors.NewValidationError("participantId", "must not be empty")
	}
	if strings.TrimSpace(r.Content) == "" {
		return errors.NewValidationError("content", "must not be empty")
	}
	return nil
}
