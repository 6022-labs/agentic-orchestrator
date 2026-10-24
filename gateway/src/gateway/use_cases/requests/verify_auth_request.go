package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
)

type VerifyAuthRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

func (r *VerifyAuthRequest) ValidateAndSanitize() error {
	r.Address = strings.TrimSpace(r.Address)
	r.Signature = strings.TrimSpace(r.Signature)

	if !helpers.IsEthAddress(r.Address) {
		return errors.NewValidationError("address", "must be a 0x-prefixed 20-byte hex address")
	}
	r.Address = strings.ToLower(r.Address)
	if !helpers.IsEthSignature(r.Signature) {
		return errors.NewValidationError("signature", "must be a 0x-prefixed 65-byte hex signature")
	}
	return nil
}
