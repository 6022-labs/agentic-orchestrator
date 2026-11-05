package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCheckingRefreshSessionRevocationTestingSuite struct {
	sut              *auth.RefreshSessionStore
	cacheStorageMock *interfaces_mocks.MockCacheStorageInterface
}

func WhenCheckingRefreshSessionRevocationBeforeEach(t *testing.T) *WhenCheckingRefreshSessionRevocationTestingSuite {
	cacheStorageMock := interfaces_mocks.NewMockCacheStorageInterface(gomock.NewController(t))

	return &WhenCheckingRefreshSessionRevocationTestingSuite{
		sut:              auth.NewRefreshSessionStore(cacheStorageMock),
		cacheStorageMock: cacheStorageMock,
	}
}

func TestWhenCheckingRefreshSessionRevocation(t *testing.T) {
	t.Parallel()

	cachePath := "auth/refresh/session-1/revoked"

	t.Run("Given the session was revoked", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingRefreshSessionRevocationTestingSuite) {
			marker := "1"
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(&marker, nil).Times(1)
		}

		t.Run("Should report it revoked", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingRefreshSessionRevocationBeforeEach(t)
			initSuite(suite)

			revoked, err := suite.sut.IsRevoked(context.Background(), "session-1")

			assert.Nil(t, err)
			assert.True(t, revoked)
		})
	})

	t.Run("Given the session was never revoked", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingRefreshSessionRevocationTestingSuite) {
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(nil, nil).Times(1)
		}

		t.Run("Should report it active", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingRefreshSessionRevocationBeforeEach(t)
			initSuite(suite)

			revoked, err := suite.sut.IsRevoked(context.Background(), "session-1")

			assert.Nil(t, err)
			assert.False(t, revoked)
		})
	})

	t.Run("Given the cache storage fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCheckingRefreshSessionRevocationTestingSuite) {
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(nil, errors.New("redis down")).Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingRefreshSessionRevocationBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.IsRevoked(context.Background(), "session-1")

			assert.ErrorContains(t, err, "redis down")
		})
	})
}
