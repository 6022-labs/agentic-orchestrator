package requests

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
)

var signedRequestNoncePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

type AuthenticateSignedRequestRequest struct {
	Method    string `header:"-"`
	Path      string `header:"-"`
	Query     string `header:"-"`
	Body      []byte `header:"-"`
	Timestamp string `header:"X-Agentic-Orchestrator-Timestamp"`
	Nonce     string `header:"X-Agentic-Orchestrator-Nonce"`
	Signature string `header:"X-Agentic-Orchestrator-Signature"`
}

func (r *AuthenticateSignedRequestRequest) ValidateAndSanitize() error {
	r.Timestamp = strings.TrimSpace(r.Timestamp)
	r.Nonce = strings.TrimSpace(r.Nonce)
	r.Signature = strings.TrimSpace(r.Signature)

	if timestamp, err := strconv.ParseInt(r.Timestamp, 10, 64); err != nil || timestamp <= 0 {
		return errors.NewValidationError("timestamp", "must be a positive Unix time in seconds")
	}
	if !signedRequestNoncePattern.MatchString(r.Nonce) {
		return errors.NewValidationError("nonce", "must be 16 to 128 letters, digits, '-' or '_'")
	}
	if !helpers.IsEthSignature(r.Signature) {
		return errors.NewValidationError("signature", "must be a 0x-prefixed 65-byte hex signature")
	}
	return nil
}
