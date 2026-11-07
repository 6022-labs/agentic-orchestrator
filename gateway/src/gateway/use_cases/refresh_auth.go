package use_cases

import (
	"context"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

// refreshReuseGrace lets tabs racing on the same refresh token share one successor instead of tripping reuse detection.
const refreshReuseGrace = 30 * time.Second

type RefreshAuthInterface interface {
	Execute(ctx context.Context, request *requests.RefreshTokenRequest) (*responses.RefreshAuthResponse, error)
}

type RefreshAuth struct {
	refreshTokenManager interfaces.RefreshTokenManagerInterface
	refreshSessionStore interfaces.RefreshSessionStoreInterface
	jwtManager          interfaces.JwtManagerInterface
	now                 func() time.Time
}

func NewRefreshAuth(
	refreshTokenManager interfaces.RefreshTokenManagerInterface,
	refreshSessionStore interfaces.RefreshSessionStoreInterface,
	jwtManager interfaces.JwtManagerInterface,
) *RefreshAuth {
	return &RefreshAuth{
		refreshTokenManager: refreshTokenManager,
		refreshSessionStore: refreshSessionStore,
		jwtManager:          jwtManager,
		now:                 time.Now,
	}
}

func (uc *RefreshAuth) Execute(
	ctx context.Context,
	request *requests.RefreshTokenRequest,
) (*responses.RefreshAuthResponse, error) {
	if strings.TrimSpace(request.RefreshToken) == "" {
		return nil, errors.NewUnauthorizedError("refresh_token_missing", "refresh token is missing")
	}

	claims, err := uc.refreshTokenManager.Parse(request.RefreshToken)
	if err != nil {
		return nil, errors.NewUnauthorizedError("refresh_token_invalid", "refresh token is invalid or expired")
	}

	revoked, err := uc.refreshSessionStore.IsRevoked(ctx, claims.SessionId)
	if err != nil {
		return nil, errors.NewInternalError("refresh_failed", err)
	}
	if revoked {
		return nil, errors.NewUnauthorizedError("session_revoked", "session was revoked")
	}

	successorId, err := uc.refreshTokenManager.NewTokenId()
	if err != nil {
		return nil, errors.NewInternalError("refresh_failed", err)
	}

	now := uc.now()
	spent := interfaces.SpentRefreshToken{SuccessorId: successorId, SpentAt: now}
	prior, spentNow, err := uc.refreshSessionStore.Spend(
		ctx,
		claims.SessionId,
		claims.TokenId,
		spent,
		claims.SessionExpiresAt,
	)
	if err != nil {
		return nil, errors.NewInternalError("refresh_failed", err)
	}
	if !spentNow {
		if prior == nil || now.Sub(prior.SpentAt) > refreshReuseGrace {
			if err := uc.refreshSessionStore.Revoke(ctx, claims.SessionId, claims.SessionExpiresAt); err != nil {
				return nil, errors.NewInternalError("refresh_failed", err)
			}
			return nil, errors.NewUnauthorizedError("refresh_token_reused", "refresh token was already used")
		}
		successorId = prior.SuccessorId
	}

	token, expiresAt, err := uc.jwtManager.CreateJwtToken(claims.Address)
	if err != nil {
		return nil, errors.NewInternalError("refresh_failed", err)
	}

	refreshToken, refreshClaims, err := uc.refreshTokenManager.Reissue(*claims, successorId)
	if err != nil {
		return nil, errors.NewInternalError("refresh_failed", err)
	}

	return &responses.RefreshAuthResponse{
		Token:            token,
		ExpiresAt:        expiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshClaims.ExpiresAt,
	}, nil
}
