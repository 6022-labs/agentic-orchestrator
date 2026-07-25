package requests

import (
	"strconv"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/google/uuid"
)

type RegisterConversationRequest struct {
	BridgeID       uuid.UUID   `json:"bridgeId"`
	Description    *string     `json:"description,omitempty"`
	ParticipantIDs []uuid.UUID `json:"participantIds,omitempty"`
}

func (r *RegisterConversationRequest) Validate() error {
	if r.BridgeID == uuid.Nil {
		return errors.NewValidationError("bridgeId", "must not be empty")
	}
	for i, pid := range r.ParticipantIDs {
		if pid == uuid.Nil {
			return errors.NewValidationError("participantIds", "entry at index "+strconv.Itoa(i)+" must not be empty")
		}
	}
	return nil
}
