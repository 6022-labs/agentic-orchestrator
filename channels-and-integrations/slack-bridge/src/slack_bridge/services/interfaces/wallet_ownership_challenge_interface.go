package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type WalletOwnershipChallengeInterface interface {
	IssueChallenge(ctx context.Context, userId string) (*uuid.UUID, error)
	ConsumeChallenge(ctx context.Context, nonce *uuid.UUID) (*string, error)
}
