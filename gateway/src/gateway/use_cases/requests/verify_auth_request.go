package requests

import (
	"regexp"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
)

var (
	ethAddressPattern   = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	ethSignaturePattern = regexp.MustCompile(`^0x[0-9a-fA-F]{130}$`)
)

type VerifyAuthRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
}

func (r *VerifyAuthRequest) Validate() error {
	r.Address = strings.TrimSpace(r.Address)
	r.Signature = strings.TrimSpace(r.Signature)
	r.Nonce = strings.TrimSpace(r.Nonce)

	if !ethAddressPattern.MatchString(r.Address) {
		return errors.NewValidationError("address", "must be a 0x-prefixed 20-byte hex address")
	}
	if !ethSignaturePattern.MatchString(r.Signature) {
		return errors.NewValidationError("signature", "must be a 0x-prefixed 65-byte hex signature")
	}
	if len(r.Nonce) == 0 {
		return errors.NewValidationError("nonce", "cannot be empty")
	}
	return nil
}
