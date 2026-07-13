package auth_test

import (
	"errors"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenConsumingNonceTestSuite struct {
	sut              *auth.NonceStore
	cacheStorageMock *interfaces_mocks.MockCacheStorageInterface
}

func WhenConsumingNonceBeforeEach(t *testing.T) *WhenConsumingNonceTestSuite {
	mockController := gomock.NewController(t)
	cacheStorageMock := interfaces_mocks.NewMockCacheStorageInterface(mockController)

	return &WhenConsumingNonceTestSuite{
		sut:              auth.NewNonceStore(cacheStorageMock),
		cacheStorageMock: cacheStorageMock,
	}
}

func TestWhenConsumingNonce(t *testing.T) {
	t.Parallel()

	t.Run("Given nonce exists", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenConsumingNonceTestSuite {
			suite := WhenConsumingNonceBeforeEach(t)
			suite.cacheStorageMock.EXPECT().ReadAndDelete("auth/nonces/abc").Return("1", nil).Times(1)
			return suite
		}

		t.Run("Should consume it exactly once", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			assert.True(t, suite.sut.Consume("abc"))
		})
	})

	t.Run("Given nonce is unknown or already used", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenConsumingNonceTestSuite {
			suite := WhenConsumingNonceBeforeEach(t)
			suite.cacheStorageMock.EXPECT().ReadAndDelete("auth/nonces/abc").Return("", errors.New("redis: nil")).Times(1)
			return suite
		}

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			assert.False(t, suite.sut.Consume("abc"))
		})
	})
}
