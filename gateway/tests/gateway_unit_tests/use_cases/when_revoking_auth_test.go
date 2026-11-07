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

type WhenRevokingAuthTestingSuite struct {
	sut *use_cases.RevokeAuth

	mockRefreshTokenManager *interfaces_mocks.MockRefreshTokenManagerInterface
	mockRefreshSessionStore *interfaces_mocks.MockRefreshSessionStoreInterface
}

func WhenRevokingAuthBeforeEach(t *testing.T) *WhenRevokingAuthTestingSuite {
	mockController := gomock.NewController(t)

	mockRefreshTokenManager := interfaces_mocks.NewMockRefreshTokenManagerInterface(mockController)
	mockRefreshSessionStore := interfaces_mocks.NewMockRefreshSessionStoreInterface(mockController)

	return &WhenRevokingAuthTestingSuite{
		sut: use_cases.NewRevokeAuth(mockRefreshTokenManager, mockRefreshSessionStore),

		mockRefreshTokenManager: mockRefreshTokenManager,
		mockRefreshSessionStore: mockRefreshSessionStore,
	}
}

func TestWhenRevokingAuth(t *testing.T) {
	t.Parallel()

	claims := &interfaces.RefreshClaims{
		Address:          common_tests.RandomEthAddress().Hex(),
		SessionId:        "session-1",
		TokenId:          "token-1",
		ExpiresAt:        time.Now().Add(time.Hour),
		SessionExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	t.Run("Given no refresh token", func(t *testing.T) {
		t.Parallel()

		t.Run("Should succeed without revoking anything", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingAuthBeforeEach(t)

			err := suite.sut.Execute(context.Background(), &requests.RefreshTokenRequest{})

			assert.NoError(t, err)
		})
	})

	t.Run("Given the refresh token is invalid or expired", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRevokingAuthTestingSuite) {
			suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(nil, errors.New("expired")).Times(1)
		}

		t.Run("Should succeed without revoking anything", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingAuthBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Execute(context.Background(), &requests.RefreshTokenRequest{RefreshToken: "refresh-token"})

			assert.NoError(t, err)
		})
	})

	t.Run("Given the session store fails", func(t *testing.T) {
		t.Parallel()

		depErr := errors.New("redis down")

		initSuite := func(suite *WhenRevokingAuthTestingSuite) {
			suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(claims, nil).Times(1)
			suite.mockRefreshSessionStore.EXPECT().Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(depErr).Times(1)
		}

		t.Run("Should return an internal error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingAuthBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Execute(context.Background(), &requests.RefreshTokenRequest{RefreshToken: "refresh-token"})

			var internalErr *cb_errors.InternalError
			assert.ErrorAs(t, err, &internalErr)
			assert.Equal(t, "revoke_failed", internalErr.Code)
			assert.Equal(t, cb_errors.InternalErrorMessage, internalErr.Message)
			assert.ErrorIs(t, internalErr, depErr)
		})
	})

	t.Run("Given a valid refresh token", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRevokingAuthTestingSuite) {
			suite.mockRefreshTokenManager.EXPECT().Parse("refresh-token").Return(claims, nil).Times(1)
			suite.mockRefreshSessionStore.EXPECT().
				Revoke(gomock.Any(), "session-1", claims.SessionExpiresAt).
				Return(nil).
				Times(1)
		}

		t.Run("Should revoke its session until it expires", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingAuthBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Execute(context.Background(), &requests.RefreshTokenRequest{RefreshToken: "refresh-token"})

			assert.NoError(t, err)
		})
	})
}
