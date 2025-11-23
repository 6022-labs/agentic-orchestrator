package requests

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type SetWalletOwnershipToSlackUserRequest struct {
	Challenge     uuid.UUID `json:"challenge"`
	WalletAddress string    `json:"walletAddress"`
}

func (r *SetWalletOwnershipToSlackUserRequest) Validate() error {
	if r.Challenge == uuid.Nil {
		return errors.New("challenge is required")
	}

	if !common.IsHexAddress(r.WalletAddress) {
		return errors.New("walletAddress must be a valid hex address")
	}

	return nil
}
