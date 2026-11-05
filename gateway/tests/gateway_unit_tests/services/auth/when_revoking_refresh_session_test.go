package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenRevokingRefreshSessionTestingSuite struct {
	sut              *auth.RefreshSessionStore
	cacheStorageMock *interfaces_mocks.MockCacheStorageInterface
}

func WhenRevokingRefreshSessionBeforeEach(t *testing.T) *WhenRevokingRefreshSessionTestingSuite {
	cacheStorageMock := interfaces_mocks.NewMockCacheStorageInterface(gomock.NewController(t))

	return &WhenRevokingRefreshSessionTestingSuite{
		sut:              auth.NewRefreshSessionStore(cacheStorageMock),
		cacheStorageMock: cacheStorageMock,
	}
}

func TestWhenRevokingRefreshSession(t *testing.T) {
	t.Parallel()

	cachePath := "auth/refresh/session-1/revoked"

	t.Run("Given the cache storage fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRevokingRefreshSessionTestingSuite) {
			suite.cacheStorageMock.EXPECT().Store(gomock.Any(), cachePath, "1", gomock.Any()).
				Return(errors.New("redis down")).Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingRefreshSessionBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Revoke(context.Background(), "session-1", time.Now().Add(time.Hour))

			assert.ErrorContains(t, err, "redis down")
		})
	})
	t.Run("Given the session is still running", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRevokingRefreshSessionTestingSuite) {
			suite.cacheStorageMock.EXPECT().
				Store(gomock.Any(), cachePath, "1", gomock.Cond(func(x any) bool {
					ttl := *x.(*time.Duration)
					return ttl > 59*time.Minute && ttl <= time.Hour
				})).
				Return(nil).
				Times(1)
		}

		t.Run("Should mark it revoked until it expires", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingRefreshSessionBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Revoke(context.Background(), "session-1", time.Now().Add(time.Hour))

			assert.Nil(t, err)
		})
	})

	t.Run("Given the session already expired", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenRevokingRefreshSessionTestingSuite) {
			suite.cacheStorageMock.EXPECT().
				Store(gomock.Any(), cachePath, "1", gomock.Cond(func(x any) bool {
					return *x.(*time.Duration) == time.Second
				})).
				Return(nil).
				Times(1)
		}

		t.Run("Should still store the marker with a one second ttl", func(t *testing.T) {
			t.Parallel()

			suite := WhenRevokingRefreshSessionBeforeEach(t)
			initSuite(suite)

			err := suite.sut.Revoke(context.Background(), "session-1", time.Now().Add(-time.Hour))

			assert.Nil(t, err)
		})
	})
}
