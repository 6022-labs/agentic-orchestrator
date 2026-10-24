package requests

import (
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type SetWalletOwnershipToSlackUserRequest struct {
	Challenge     uuid.UUID `json:"challenge"`
	WalletAddress string    `json:"walletAddress"`
}

// ValidateAndSanitize lowercases the wallet address so wallet comparisons are case-insensitive.
func (r *SetWalletOwnershipToSlackUserRequest) ValidateAndSanitize() error {
	if r.Challenge == uuid.Nil {
		return common_errors.NewValidationError("challenge", "is required")
	}

	if !strings.HasPrefix(r.WalletAddress, "0x") || !common.IsHexAddress(r.WalletAddress) {
		return common_errors.NewValidationError("walletAddress", "must be a 0x-prefixed hex address")
	}

	r.WalletAddress = strings.ToLower(r.WalletAddress)
	return nil
}
