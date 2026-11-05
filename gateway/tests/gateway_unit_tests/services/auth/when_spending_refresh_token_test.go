package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenSpendingRefreshTokenTestingSuite struct {
	sut              *auth.RefreshSessionStore
	cacheStorageMock *interfaces_mocks.MockCacheStorageInterface
}

func WhenSpendingRefreshTokenBeforeEach(t *testing.T) *WhenSpendingRefreshTokenTestingSuite {
	cacheStorageMock := interfaces_mocks.NewMockCacheStorageInterface(gomock.NewController(t))

	return &WhenSpendingRefreshTokenTestingSuite{
		sut:              auth.NewRefreshSessionStore(cacheStorageMock),
		cacheStorageMock: cacheStorageMock,
	}
}

func TestWhenSpendingRefreshToken(t *testing.T) {
	t.Parallel()

	cachePath := "auth/refresh/session-1/spent/token-1"
	spent := interfaces.SpentRefreshToken{SuccessorId: "token-2", SpentAt: time.Now().UTC().Truncate(time.Second)}
	payload, _ := json.Marshal(spent)
	expiresAt := time.Now().Add(time.Hour)
	ttlUntilExpiry := gomock.Cond(func(x any) bool {
		ttl := x.(time.Duration)
		return ttl > 59*time.Minute && ttl <= time.Hour
	})
	expectAlreadyStored := func(suite *WhenSpendingRefreshTokenTestingSuite) {
		suite.cacheStorageMock.EXPECT().StoreIfAbsent(gomock.Any(), cachePath, string(payload), gomock.Any()).
			Return(false, nil).Times(1)
	}

	t.Run("Given the cache storage fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSpendingRefreshTokenTestingSuite) {
			suite.cacheStorageMock.EXPECT().StoreIfAbsent(gomock.Any(), cachePath, string(payload), gomock.Any()).
				Return(false, errors.New("redis down")).Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSpendingRefreshTokenBeforeEach(t)
			initSuite(suite)

			_, spentNow, err := suite.sut.Spend(context.Background(), "session-1", "token-1", spent, expiresAt)

			assert.ErrorContains(t, err, "redis down")
			assert.False(t, spentNow)
		})
	})

	t.Run("Given reading the spend record fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSpendingRefreshTokenTestingSuite) {
			expectAlreadyStored(suite)
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(nil, errors.New("redis down")).Times(1)
		}

		t.Run("Should return the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSpendingRefreshTokenBeforeEach(t)
			initSuite(suite)

			_, spentNow, err := suite.sut.Spend(context.Background(), "session-1", "token-1", spent, expiresAt)

			assert.ErrorContains(t, err, "redis down")
			assert.False(t, spentNow)
		})
	})

	t.Run("Given the spend record expired between the two calls", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSpendingRefreshTokenTestingSuite) {
			expectAlreadyStored(suite)
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(nil, nil).Times(1)
		}

		t.Run("Should report it as already spent without a record", func(t *testing.T) {
			t.Parallel()

			suite := WhenSpendingRefreshTokenBeforeEach(t)
			initSuite(suite)

			prior, spentNow, err := suite.sut.Spend(context.Background(), "session-1", "token-1", spent, expiresAt)

			assert.Nil(t, err)
			assert.False(t, spentNow)
			assert.Nil(t, prior)
		})
	})

	t.Run("Given the spend record is corrupt", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSpendingRefreshTokenTestingSuite) {
			corrupt := "not-json"
			expectAlreadyStored(suite)
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(&corrupt, nil).Times(1)
		}

		t.Run("Should return a decode error", func(t *testing.T) {
			t.Parallel()

			suite := WhenSpendingRefreshTokenBeforeEach(t)
			initSuite(suite)

			_, spentNow, err := suite.sut.Spend(context.Background(), "session-1", "token-1", spent, expiresAt)

			assert.ErrorContains(t, err, "failed to decode spent refresh token")
			assert.False(t, spentNow)
		})
	})

	t.Run("Given the token was never spent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenSpendingRefreshTokenTestingSuite) {
			suite.cacheStorageMock.EXPECT().
				StoreIfAbsent(gomock.Any(), cachePath, string(payload), ttlUntilExpiry).
				Return(true, nil).
				Times(1)
		}

		t.Run("Should spend it until the session expires", func(t *testing.T) {
			t.Parallel()

			suite := WhenSpendingRefreshTokenBeforeEach(t)
			initSuite(suite)

			prior, spentNow, err := suite.sut.Spend(context.Background(), "session-1", "token-1", spent, expiresAt)

			assert.Nil(t, err)
			assert.True(t, spentNow)
			assert.Equal(t, &spent, prior)
		})
	})

	t.Run("Given the token was already spent", func(t *testing.T) {
		t.Parallel()

		first := interfaces.SpentRefreshToken{SuccessorId: "token-first", SpentAt: spent.SpentAt.Add(-time.Second)}
		firstPayload, _ := json.Marshal(first)

		initSuite := func(suite *WhenSpendingRefreshTokenTestingSuite) {
			stringPayload := string(firstPayload)
			expectAlreadyStored(suite)
			suite.cacheStorageMock.EXPECT().Read(gomock.Any(), cachePath).Return(&stringPayload, nil).Times(1)
		}

		t.Run("Should return the first spend record", func(t *testing.T) {
			t.Parallel()

			suite := WhenSpendingRefreshTokenBeforeEach(t)
			initSuite(suite)

			prior, spentNow, err := suite.sut.Spend(context.Background(), "session-1", "token-1", spent, expiresAt)

			assert.Nil(t, err)
			assert.False(t, spentNow)
			assert.Equal(t, first.SuccessorId, prior.SuccessorId)
			assert.True(t, first.SpentAt.Equal(prior.SpentAt))
		})
	})
}
