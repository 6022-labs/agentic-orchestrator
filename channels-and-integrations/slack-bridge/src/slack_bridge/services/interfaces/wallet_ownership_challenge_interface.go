package interfaces

import (
	"github.com/google/uuid"
)

type WalletOwnershipChallengeInterface interface {
	IssueChallenge(userId string) (*uuid.UUID, error)
	ConsumeChallenge(nonce *uuid.UUID) (*string, error)
}
