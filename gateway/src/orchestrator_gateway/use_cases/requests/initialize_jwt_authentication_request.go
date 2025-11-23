package requests

import (
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/ethereum/go-ethereum/common"
)

type InitializeJwtAuthenticationRequest struct {
	Address string `json:"address"`
}

func (r *InitializeJwtAuthenticationRequest) Validate() error {
	if len(strings.TrimSpace(r.Address)) == 0 {
		return errors.NewValidationError("address", "Address is required")
	}

	if !common.IsHexAddress(r.Address) {
		return errors.NewValidationError("address", "Invalid ethereum address")
	}

	return nil
}
