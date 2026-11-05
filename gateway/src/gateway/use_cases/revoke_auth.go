package use_cases

import (
	"context"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type RevokeAuthInterface interface {
	// Execute is idempotent: a missing or unreadable token has no session left to revoke.
	Execute(ctx context.Context, request *requests.RefreshTokenRequest) error
}

type RevokeAuth struct {
	refreshTokenManager interfaces.RefreshTokenManagerInterface
	refreshSessionStore interfaces.RefreshSessionStoreInterface
}

func NewRevokeAuth(
	refreshTokenManager interfaces.RefreshTokenManagerInterface,
	refreshSessionStore interfaces.RefreshSessionStoreInterface,
) *RevokeAuth {
	return &RevokeAuth{
		refreshTokenManager: refreshTokenManager,
		refreshSessionStore: refreshSessionStore,
	}
}

func (uc *RevokeAuth) Execute(ctx context.Context, request *requests.RefreshTokenRequest) error {
	if strings.TrimSpace(request.RefreshToken) == "" {
		return nil
	}

	claims, err := uc.refreshTokenManager.Parse(request.RefreshToken)
	if err != nil {
		return nil
	}

	if err := uc.refreshSessionStore.Revoke(ctx, claims.SessionId, claims.SessionExpiresAt); err != nil {
		return errors.NewInternalError("revoke_failed", err.Error())
	}
	return nil
}
