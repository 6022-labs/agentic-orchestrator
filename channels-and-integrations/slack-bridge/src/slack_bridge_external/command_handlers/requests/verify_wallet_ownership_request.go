package requests

import (
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
)

type VerifyWalletOwnershipRequest struct {
	UserId string `form:"userId"`
}

func (r *VerifyWalletOwnershipRequest) Validate() error {
	if len(strings.TrimSpace(r.UserId)) == 0 {
		return common_errors.NewValidationError("userId", "cannot be empty")
	}

	return nil
}
