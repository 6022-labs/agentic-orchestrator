package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
)

type FinishJwtAuthenticationRequest struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (r *FinishJwtAuthenticationRequest) Validate() error {
	if len(strings.TrimSpace(r.Message)) == 0 {
		return errors.NewValidationError("message", "is required")
	}

	if len(strings.TrimSpace(r.Signature)) == 0 {
		return errors.NewValidationError("signature", "is required")
	}

	return nil
}
