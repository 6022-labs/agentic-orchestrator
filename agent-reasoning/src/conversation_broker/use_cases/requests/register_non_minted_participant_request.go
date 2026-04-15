package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
)

type RegisterNonMintedParticipantRequest struct {
	DisplayName string `json:"displayName"`
}

func (r *RegisterNonMintedParticipantRequest) Validate() error {
	if strings.TrimSpace(r.DisplayName) == "" {
		return errors.NewValidationError("displayName", "must not be empty")
	}
	return nil
}
