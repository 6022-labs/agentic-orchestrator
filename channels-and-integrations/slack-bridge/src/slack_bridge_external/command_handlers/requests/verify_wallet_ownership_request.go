package requests

import (
	"errors"
	"strings"
)

type VerifyWalletOwnershipRequest struct {
	UserId string `form:"userId"`
}

func (r *VerifyWalletOwnershipRequest) Validate() error {
	if len(strings.TrimSpace(r.UserId)) == 0 {
		return errors.New("user id cannot be empty")
	}

	return nil
}
