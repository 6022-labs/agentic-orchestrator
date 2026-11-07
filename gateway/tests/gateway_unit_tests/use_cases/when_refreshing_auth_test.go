package use_cases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRefreshingAuthTestingSuite struct {
	sut *use_cases.RefreshAuth

	mockRefreshTokenManager *interfaces_mocks.MockRefreshTokenManagerInterface
	mockRefreshSessionStore *interfaces_mocks.MockRefreshSessionStoreInterface
	mockJwtManager          *interfaces_mocks.MockJwtManagerInterface
}

func WhenRefreshingAuthBeforeEach(t *testing.T) *WhenRefreshingAuthTestingSuite {
	mockController := gomock.NewController(t)

	mockRefreshTokenManager := interfaces_mocks.NewMockRefreshTokenManagerInterface(mockController)
	mockRefreshSessionStore := interfaces_mocks.NewMockRefreshSessionStoreInterface(mockController)
	mockJwtManager := interfaces_mocks.NewMockJwtManagerInterface(mockController)

	return &WhenRefreshingAuthTestingSuite{
		sut: use_cases.NewRefreshAuth(mockRefreshTokenManager, mockRefreshSessionStore, mockJwtManager),

		mockRefreshTokenManager: mockRefreshTokenManager,
		mockRefreshSessionStore: mockRefreshSessionStore,
		mockJwtManager:          mockJwtManager,
	}
}

func TestWhenRefreshingAuth(t *testing.T) {
	t.Parallel()

	claims := &interfaces.RefreshClaims{
		Address:          common_tests.RandomEthAddress().Hex(),
		SessionId:        "session-1",
		TokenId:          "token-1",
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		SessionExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	accessExpiresAt := time.Now().Add(time.Hour)
	request := &requests.RefreshTokenRequest{RefreshToken: "refresh-token"}
	expectLiveSession := func(suite *WhenRefreshingAuthTestingSuite) {
		suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(claims, nil).Times(1)
		suite.mockRefreshSessionStore.EXPECT().IsRevoked(gomock.Any(), "session-1").Return(false, nil).Times(1)
	}
	expectSuccessorId := func(suite *WhenRefreshingAuthTestingSuite) {
		expectLiveSession(suite)
		suite.mockRefreshTokenManager.EXPECT().NewTokenId().Return("token-2", nil).Times(1)
	}
	expectSpent := func(suite *WhenRefreshingAuthTestingSuite) {
		suite.mockRefreshSessionStore.EXPECT().
			Spend(gomock.Any(), "session-1", "token-1", gomock.Any(), claims.SessionExpiresAt).
			Return(nil, true, nil).
			Times(1)
	}
	expectSpentOutsideGrace := func(suite *WhenRefreshingAuthTestingSuite) {
		suite.mockRefreshSessionStore.EXPECT().
			Spend(gomock.Any(), "session-1", "token-1", gomock.Any(), claims.SessionExpiresAt).
			Return(&interfaces.SpentRefreshToken{SuccessorId: "token-old", SpentAt: time.Now().Add(-time.Minute)}, false, nil).
			Times(1)
	}
	expectIssued := func(suite *WhenRefreshingAuthTestingSuite, successorId string) {
		suite.mockJwtManager.EXPECT().CreateJwtToken(claims.Address).Return("jwt-token", accessExpiresAt, nil).Times(1)
		suite.mockRefreshTokenManager.EXPECT().
			Reissue(*claims, successorId).
			Return("next-refresh-token", &interfaces.RefreshClaims{ExpiresAt: claims.ExpiresAt}, nil).
			Times(1)
	}

	t.Run("Given no refresh token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)

			_, err := suite.sut.Execute(context.Background(), &requests.RefreshTokenRequest{})

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "refresh_token_missing", unauthorizedErr.Code)
		})
	})

	t.Run("Given the refresh token is invalid or expired", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(nil, errors.New("expired")).Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "refresh_token_invalid", unauthorizedErr.Code)
		})
	})

	t.Run("Given the revocation check fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("redis down")

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(claims, nil).Times(1)
			suite.mockRefreshSessionStore.EXPECT().IsRevoked(gomock.Any(), "session-1").
				Return(false, depErr).Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "refresh_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given the session was revoked", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(claims, nil).Times(1)
			suite.mockRefreshSessionStore.EXPECT().IsRevoked(gomock.Any(), "session-1").Return(true, nil).Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "session_revoked", unauthorizedErr.Code)
		})
	})

	t.Run("Given the token id generation fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("entropy exhausted")

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectLiveSession(suite)
			suite.mockRefreshTokenManager.EXPECT().NewTokenId().Return("", depErr).Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "refresh_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given spending the refresh token fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("redis down")

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Spend(gomock.Any(), "session-1", "token-1", gomock.Any(), claims.SessionExpiresAt).
				Return(nil, false, depErr).
				Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "refresh_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given the session revocation fails after a reuse", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("redis down")

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			expectSpentOutsideGrace(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(depErr).
				Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "refresh_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given the refresh token was already used outside the grace window", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			expectSpentOutsideGrace(suite)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(nil).
				Times(1)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "refresh_token_reused", unauthorizedErr.Code)
		})

		t.Run("Should revoke the session", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(nil).
				Times(1)

			_, _ = suite.sut.Execute(context.Background(), request)
		})
	})

	t.Run("Given the refresh token was already used and its record is gone", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Spend(gomock.Any(), "session-1", "token-1", gomock.Any(), claims.SessionExpiresAt).
				Return(nil, false, nil).
				Times(1)
		}

		t.Run("Should return an unauthorized error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(nil).
				Times(1)

			_, err := suite.sut.Execute(context.Background(), request)

			var unauthorizedErr *cb_errors.UnauthorizedError
			assert.ErrorAs(t, err, &unauthorizedErr)
			assert.Equal(t, "refresh_token_reused", unauthorizedErr.Code)
		})

		t.Run("Should revoke the session", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(nil).
				Times(1)

			_, _ = suite.sut.Execute(context.Background(), request)
		})
	})

	t.Run("Given a concurrent refresh spent the token within the grace window", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Spend(gomock.Any(), "session-1", "token-1", gomock.Any(), claims.SessionExpiresAt).
				Return(&interfaces.SpentRefreshToken{SuccessorId: "token-shared", SpentAt: time.Now()}, false, nil).
				Times(1)
			expectIssued(suite, "token-shared")
		}

		t.Run("Should reissue the concurrent refresh's successor", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "next-refresh-token", response.RefreshToken)
		})
	})

	t.Run("Given the access token creation fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("sign failed")

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			expectSpent(suite)
			suite.mockJwtManager.EXPECT().
				CreateJwtToken(claims.Address).
				Return("", time.Time{}, depErr).
				Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "refresh_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given the refresh token reissue fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("sign failed")

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			expectSpent(suite)
			suite.mockJwtManager.EXPECT().
				CreateJwtToken(claims.Address).
				Return("jwt-token", accessExpiresAt, nil).
				Times(1)
			suite.mockRefreshTokenManager.EXPECT().
				Reissue(*claims, "token-2").
				Return("", nil, depErr).
				Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.Execute(context.Background(), request)

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "refresh_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given a valid unused refresh token", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRefreshingAuthTestingSuite) {
			expectSuccessorId(suite)
			suite.mockRefreshSessionStore.EXPECT().
				Spend(gomock.Any(), "session-1", "token-1", gomock.Cond(func(x any) bool {
					return x.(interfaces.SpentRefreshToken).SuccessorId == "token-2"
				}), claims.SessionExpiresAt).
				Return(nil, true, nil).
				Times(1)
			expectIssued(suite, "token-2")
		}

		t.Run("Should return a fresh access token", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "jwt-token", response.Token)
			assert.Equal(t, accessExpiresAt, response.ExpiresAt)
		})

		t.Run("Should return the rotated refresh token", func(t *testing.T) {
			t.Parallel()

			suite := WhenRefreshingAuthBeforeEach(t)
			initSuite(suite)

			response, err := suite.sut.Execute(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, "next-refresh-token", response.RefreshToken)
			assert.Equal(t, claims.ExpiresAt, response.RefreshExpiresAt)
		})
	})
}
