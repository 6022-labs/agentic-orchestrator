package requests

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type SetWalletOwnershipToSlackUserRequest struct {
	Challenge     uuid.UUID `json:"challenge"`
	WalletAddress string    `json:"walletAddress"`
}

func (r *SetWalletOwnershipToSlackUserRequest) Validate() error {
	if r.Challenge == uuid.Nil {
		return fmt.Errorf("%w: challenge is required", ErrValidation)
	}

	if !common.IsHexAddress(r.WalletAddress) {
		return fmt.Errorf("%w: walletAddress must be a valid hex address", ErrValidation)
	}

	return nil
}
