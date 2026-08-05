package services_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenStoringX402PaymentChallengeTestingSuite struct {
	sut *services.X402PaymentChallengeCache

	cacheStorage *interfaces_mocks.MockCacheStorageInterface
}

func WhenStoringX402PaymentChallengeBeforeEach(t *testing.T) *WhenStoringX402PaymentChallengeTestingSuite {
	cacheStorage := interfaces_mocks.NewMockCacheStorageInterface(gomock.NewController(t))

	return &WhenStoringX402PaymentChallengeTestingSuite{
		sut: services.NewX402PaymentChallengeCache(zap.NewNop(), cacheStorage, &settings.X402PaymentSettings{
			ChallengeCacheTtlSeconds: 120,
		}),
		cacheStorage: cacheStorage,
	}
}

func TestWhenStoringX402PaymentChallenge(t *testing.T) {
	t.Parallel()

	t.Run("Given a requirements document", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store it serialized under the normalized domain key with the configured TTL", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringX402PaymentChallengeBeforeEach(t)
			suite.cacheStorage.EXPECT().
				Store("x402-payment-challenges/agent.eth", gomock.Any(), gomock.Any()).
				DoAndReturn(func(key, value string, expiration *time.Duration) error {
					assert.Contains(t, value, `"x402Version":2`)
					assert.Equal(t, 120*time.Second, *expiration)
					return nil
				})

			suite.sut.Store(" Agent.ETH ", &models.X402PaymentChallenge{X402Version: models.X402Version})
		})
	})

	t.Run("Given the cache storage fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should swallow the error", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringX402PaymentChallengeBeforeEach(t)
			suite.cacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(assert.AnError)

			suite.sut.Store("agent.eth", &models.X402PaymentChallenge{X402Version: models.X402Version})
		})
	})
}
