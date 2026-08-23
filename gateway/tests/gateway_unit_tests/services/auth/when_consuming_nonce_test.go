package auth_test

import (
	"context"

	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/auth"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	metrics_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/metrics_mocks/interfaces_mocks"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenConsumingNonceTestingSuite struct {
	gatewayMetricsMock *metrics_mocks.MockGatewayMetricsInterface
	sut                *auth.NonceStore
	cacheStorageMock   *interfaces_mocks.MockCacheStorageInterface
}

func WhenConsumingNonceBeforeEach(t *testing.T) *WhenConsumingNonceTestingSuite {
	mockController := gomock.NewController(t)
	cacheStorageMock := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	gatewayMetricsMock := metrics_mocks.NewMockGatewayMetricsInterface(mockController)

	return &WhenConsumingNonceTestingSuite{
		sut: auth.NewNonceStore(
			cacheStorageMock,
			&settings.JwtSettings{NonceTtl: 300},
			gatewayMetricsMock,
		),
		gatewayMetricsMock: gatewayMetricsMock,
		cacheStorageMock:   cacheStorageMock,
	}
}

func TestWhenConsumingNonce(t *testing.T) {
	t.Parallel()

	wallet := common.HexToAddress("0x4bca9E17a8326d283d1fC68A725B46b42c2D9146")
	cachePath := "auth/challenges/0x4bca9e17a8326d283d1fc68a725b46b42c2d9146"

	t.Run("Given a challenge exists for the wallet", func(t *testing.T) {
		t.Parallel()

		stored := interfaces.ChallengeEntry{
			Nonce:     "nonce-abc",
			ChainId:   80002,
			IssuedAt:  time.Now().UTC().Truncate(time.Second),
			ExpiresAt: time.Now().UTC().Truncate(time.Second).Add(5 * time.Minute),
		}
		payload, _ := json.Marshal(stored)

		initSuite := func() *WhenConsumingNonceTestingSuite {
			suite := WhenConsumingNonceBeforeEach(t)
			stringPayload := string(payload)
			suite.cacheStorageMock.EXPECT().ReadAndDelete(gomock.Any(), cachePath).Return(&stringPayload, nil).Times(1)
			suite.gatewayMetricsMock.EXPECT().RecordCacheResult(gomock.Any(), "nonce", true)
			return suite
		}

		t.Run("Should consume it exactly once and return the challenge", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			entry := suite.sut.Consume(context.Background(), wallet)

			if assert.NotNil(t, entry) {
				assert.Equal(t, stored.Nonce, entry.Nonce)
				assert.Equal(t, stored.ChainId, entry.ChainId)
				assert.True(t, stored.IssuedAt.Equal(entry.IssuedAt))
				assert.True(t, stored.ExpiresAt.Equal(entry.ExpiresAt))
			}
		})
	})

	t.Run("Given the wallet has no challenge or it was already used", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenConsumingNonceTestingSuite {
			suite := WhenConsumingNonceBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordCacheResult(gomock.Any(), "nonce", false)
			suite.cacheStorageMock.EXPECT().
				ReadAndDelete(gomock.Any(), cachePath).
				Return(nil, nil).
				Times(1)
			return suite
		}

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			assert.Nil(t, suite.sut.Consume(context.Background(), wallet))
		})
	})

	t.Run("Given the cache is unreachable", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenConsumingNonceTestingSuite {
			suite := WhenConsumingNonceBeforeEach(t)
			// An outage is not a miss, so nothing is recorded.
			suite.cacheStorageMock.EXPECT().
				ReadAndDelete(gomock.Any(), cachePath).
				Return(nil, errors.New("dial tcp: connection refused")).
				Times(1)
			return suite
		}

		t.Run("Should reject it without counting a miss", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			assert.Nil(t, suite.sut.Consume(context.Background(), wallet))
		})
	})

	t.Run("Given the stored payload is corrupted", func(t *testing.T) {
		t.Parallel()

		initSuite := func() *WhenConsumingNonceTestingSuite {
			suite := WhenConsumingNonceBeforeEach(t)
			corrupted := "not-json"
			suite.cacheStorageMock.EXPECT().ReadAndDelete(gomock.Any(), cachePath).Return(&corrupted, nil).Times(1)
			// A corrupt payload still read cleanly from the cache: a hit, then a decode failure.
			suite.gatewayMetricsMock.EXPECT().RecordCacheResult(gomock.Any(), "nonce", true)
			return suite
		}

		t.Run("Should reject it", func(t *testing.T) {
			t.Parallel()

			suite := initSuite()

			assert.Nil(t, suite.sut.Consume(context.Background(), wallet))
		})
	})
}
