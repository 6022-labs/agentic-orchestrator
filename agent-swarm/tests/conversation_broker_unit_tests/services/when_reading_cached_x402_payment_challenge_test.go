package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenReadingCachedX402PaymentChallengeTestingSuite struct {
	sut *services.X402PaymentChallengeCache

	cacheStorage *interfaces_mocks.MockCacheStorageInterface
}

func WhenReadingCachedX402PaymentChallengeBeforeEach(t *testing.T) *WhenReadingCachedX402PaymentChallengeTestingSuite {
	mockController := gomock.NewController(t)
	cacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)

	sut := services.NewX402PaymentChallengeCache(
		zap.NewNop(),
		cacheStorage,
		&settings.X402PaymentSettings{
			ChallengeCacheTtlSeconds: 300,
		},
	)

	return &WhenReadingCachedX402PaymentChallengeTestingSuite{
		sut:          sut,
		cacheStorage: cacheStorage,
	}
}

func TestWhenReadingCachedX402PaymentChallenge(t *testing.T) {
	t.Parallel()

	t.Run("Given a cached requirements document", func(t *testing.T) {
		t.Parallel()

		t.Run("Should decode it from the normalized domain key", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingCachedX402PaymentChallengeBeforeEach(t)
			suite.cacheStorage.EXPECT().
				Read(gomock.Any(), "x402-payment-challenges/agent.eth").
				Return(`{"x402Version":2,"accepts":[{"scheme":"exact","network":"eip155:80002"}]}`, nil)

			paymentChallenge, ok := suite.sut.Read(context.Background(), " Agent.ETH ")

			assert.True(t, ok)
			assert.Equal(t, models.X402Version, paymentChallenge.X402Version)
			assert.Len(t, paymentChallenge.Accepts, 1)
		})
	})

	t.Run("Given a cache miss", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report not found", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingCachedX402PaymentChallengeBeforeEach(t)
			suite.cacheStorage.EXPECT().
				Read(gomock.Any(), "x402-payment-challenges/agent.eth").
				Return("", assert.AnError)

			paymentChallenge, ok := suite.sut.Read(context.Background(), "agent.eth")

			assert.False(t, ok)
			assert.Nil(t, paymentChallenge)
		})
	})

	t.Run("Given a corrupted cache entry", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report not found", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingCachedX402PaymentChallengeBeforeEach(t)
			suite.cacheStorage.EXPECT().
				Read(gomock.Any(), "x402-payment-challenges/agent.eth").
				Return("not-json", nil)

			paymentChallenge, ok := suite.sut.Read(context.Background(), "agent.eth")

			assert.False(t, ok)
			assert.Nil(t, paymentChallenge)
		})
	})
}
