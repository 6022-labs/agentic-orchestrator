package services_test

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenResolvingRemoteAgentCardTestingSuite struct {
	sut *services.RemoteAgentCardResolver

	mockCacheStorage             *interfaces_mocks.MockCacheStorageInterface
	mockEnsUrlResolver           *interfaces_mocks.MockEnsUrlResolverInterface
	mockRemoteAgentCardRequester *interfaces_mocks.MockRemoteAgentCardRequesterInterface
}

func WhenResolvingRemoteAgentCardBeforeEach(t *testing.T) *WhenResolvingRemoteAgentCardTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockEnsUrlResolver := interfaces_mocks.NewMockEnsUrlResolverInterface(mockController)
	mockRemoteAgentCardRequester := interfaces_mocks.NewMockRemoteAgentCardRequesterInterface(mockController)

	sut := services.NewRemoteAgentCardResolver(
		zap.NewNop(),
		mockCacheStorage,
		mockEnsUrlResolver,
		mockRemoteAgentCardRequester,
		&settings.RemoteAgentResolutionSettings{CacheTtlSeconds: 300},
	)

	return &WhenResolvingRemoteAgentCardTestingSuite{
		sut: sut,

		mockCacheStorage:             mockCacheStorage,
		mockEnsUrlResolver:           mockEnsUrlResolver,
		mockRemoteAgentCardRequester: mockRemoteAgentCardRequester,
	}
}

func cardWithResponses() *models.RemoteAgentCard {
	return &models.RemoteAgentCard{
		Protocol: "6022",
		Endpoints: models.RemoteAgentCardEndpoints{
			Responses: "https://node.example/api/responses",
		},
	}
}

func TestWhenResolvingRemoteAgentCard(t *testing.T) {
	t.Parallel()

	t.Run("Given the card is not cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read("remote-agent-cards/agent.eth").Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("https://node.example", nil).Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().GetCard("https://node.example").Return(cardWithResponses(), nil).Times(1)
			suite.mockCacheStorage.EXPECT().
				Store("remote-agent-cards/agent.eth", gomock.Any(), gomock.Any()).
				DoAndReturn(func(key, value string, ttl *time.Duration) error {
					assert.Equal(t, 300*time.Second, *ttl)
					assert.Contains(t, value, "https://node.example/api/responses")
					return nil
				}).
				Times(1)
		}

		t.Run("Should resolve the card and cache it with the TTL", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			card, err := suite.sut.ResolveCard("agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
			assert.Equal(t, "https://node.example/api/responses", card.Endpoints.Responses)
		})
	})

	t.Run("Given concurrent resolutions miss the cache for the same domain", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) (chan struct{}, chan struct{}) {
			misses := make(chan struct{}, 3)
			release := make(chan struct{})

			var storeMutex sync.Mutex
			stored := ""
			suite.mockCacheStorage.EXPECT().
				Read("remote-agent-cards/agent.eth").
				DoAndReturn(func(string) (string, error) {
					storeMutex.Lock()
					defer storeMutex.Unlock()
					if stored == "" {
						misses <- struct{}{}
						return "", errors.New("miss")
					}
					return stored, nil
				}).
				Times(4)
			suite.mockEnsUrlResolver.EXPECT().
				ResolveUrl("agent.eth").
				DoAndReturn(func(string) (string, error) {
					<-release
					return "https://node.example", nil
				}).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().GetCard("https://node.example").Return(cardWithResponses(), nil).Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(key, value string, ttl *time.Duration) error {
					storeMutex.Lock()
					defer storeMutex.Unlock()
					stored = value
					return nil
				}).
				Times(1)
			return misses, release
		}

		t.Run("Should resolve once and serve the waiter from the cache", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			misses, release := initSuite(suite)

			var wg sync.WaitGroup
			results := make([]*models.RemoteAgentCard, 2)
			for i := range results {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					card, err := suite.sut.ResolveCard("agent.eth")
					assert.NoError(t, err)
					results[i] = card
				}(i)
			}

			// Hold the resolution until both outer reads plus the lock holder's
			// re-check have missed; the waiter's re-check then must hit the cache.
			<-misses
			<-misses
			<-misses
			close(release)
			wg.Wait()

			assert.Equal(t, results[0], results[1])
		})
	})

	t.Run("Given the card is cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			cached, err := json.Marshal(cardWithResponses())
			assert.NoError(t, err)
			suite.mockCacheStorage.EXPECT().Read("remote-agent-cards/agent.eth").Return(string(cached), nil).Times(1)
		}

		t.Run("Should serve the card from the cache without resolving", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			card, err := suite.sut.ResolveCard("agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
		})
	})

	t.Run("Given the cached value is corrupted", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read("remote-agent-cards/agent.eth").Return("not-json", nil).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("https://node.example", nil).Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().GetCard("https://node.example").Return(cardWithResponses(), nil).Times(1)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		}

		t.Run("Should fall back to a full resolution", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			card, err := suite.sut.ResolveCard("agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
		})
	})

	t.Run("Given the ENS resolution fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("", errors.New("boom")).Times(1)
		}

		t.Run("Should propagate the error and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveCard("agent.eth")

			assert.Error(t, err)
		})
	})

	t.Run("Given no url is published", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("", nil).Times(1)
		}

		t.Run("Should return an error and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveCard("agent.eth")

			assert.Error(t, err)
		})
	})
}
