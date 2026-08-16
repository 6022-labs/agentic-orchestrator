package services_test

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenResolvingRegistryAgentTestingSuite struct {
	sut *services.RegistryAgentResolver

	mockCacheStorage           *interfaces_mocks.MockCacheStorageInterface
	mockAgentRegistryRequester *interfaces_mocks.MockAgentRegistryRequesterInterface
}

func WhenResolvingRegistryAgentBeforeEach(t *testing.T) *WhenResolvingRegistryAgentTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockAgentRegistryRequester := interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)

	return &WhenResolvingRegistryAgentTestingSuite{
		sut: services.NewRegistryAgentResolver(
			zap.NewNop(),
			mockCacheStorage,
			mockAgentRegistryRequester,
			&settings.RegistryAgentResolutionSettings{CacheTtlSeconds: 60},
		),

		mockCacheStorage:           mockCacheStorage,
		mockAgentRegistryRequester: mockAgentRegistryRequester,
	}
}

const (
	registryAgentEnsCachePath        = "registry-agents/ens/hermes.agents6022.eth"
	registryAgentCollectionCachePath = "registry-agents/collection/0xabc/1"
)

func registryAgentForTest() *models.RegistryAgent {
	return model_builders.NewRegistryAgentBuilder().
		WithName("hermes").
		WithRole("facilitator").
		WithEnsDomain("hermes.agents6022.eth").
		WithCollection("0xABC", 1).
		WithAddress(models.EvmAddressType, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").
		Build()
}

func TestWhenResolvingRegistryAgent(t *testing.T) {
	t.Parallel()

	t.Run("Given the agent is not cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRegistryAgentTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(registryAgentEnsCachePath).Return("", errors.New("miss")).Times(2)
			suite.mockAgentRegistryRequester.EXPECT().
				GetMintedAgentByEnsDomain("hermes.agents6022.eth").
				Return(registryAgentForTest(), nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(registryAgentEnsCachePath, gomock.Any(), gomock.Any()).
				DoAndReturn(func(key, value string, ttl *time.Duration) error {
					assert.Equal(t, 60*time.Second, *ttl)
					assert.Contains(t, value, "facilitator")
					return nil
				}).
				Times(1)
			// One fetch warms the other lookup path too.
			suite.mockCacheStorage.EXPECT().
				Store(registryAgentCollectionCachePath, gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
		}

		t.Run("Should fetch it once and cache it with the TTL", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)
			initSuite(suite)

			agent, err := suite.sut.ResolveByEnsDomain("hermes.agents6022.eth")

			assert.NoError(t, err)
			assert.Equal(t, "hermes", agent.Name)
		})

		t.Run("Should normalise the domain so casing does not split the cache", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)
			initSuite(suite)

			agent, err := suite.sut.ResolveByEnsDomain("  Hermes.Agents6022.ETH  ")

			assert.NoError(t, err)
			assert.Equal(t, "hermes", agent.Name)
		})
	})

	t.Run("Given the agent is cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRegistryAgentTestingSuite) {
			cached, err := json.Marshal(registryAgentForTest())
			assert.NoError(t, err)
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return(string(cached), nil).Times(1)
		}

		t.Run("Should serve it without touching the registry", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)
			initSuite(suite)
			suite.mockAgentRegistryRequester.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Times(0)

			agent, err := suite.sut.ResolveByCollection("0xABC", 1)

			assert.NoError(t, err)
			assert.Equal(t, "hermes", agent.Name)
		})
	})

	t.Run("Given the registry does not know the agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRegistryAgentTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockAgentRegistryRequester.EXPECT().GetAgent("0xabc", uint64(1)).Return(nil, nil).Times(1)
		}

		t.Run("Should not cache the miss, so a freshly minted agent is visible at once", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)
			initSuite(suite)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			agent, err := suite.sut.ResolveByCollection("0xabc", 1)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given the registry fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRegistryAgentTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)
			initSuite(suite)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			_, err := suite.sut.ResolveByCollection("0xabc", 1)

			assert.Error(t, err)
		})
	})

	t.Run("Given the cached value is corrupted", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRegistryAgentTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("not-json", nil).Times(2)
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				Return(registryAgentForTest(), nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
		}

		t.Run("Should fall back to a full resolution", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)
			initSuite(suite)

			agent, err := suite.sut.ResolveByCollection("0xabc", 1)

			assert.NoError(t, err)
			assert.Equal(t, "hermes", agent.Name)
		})
	})

	t.Run("Given concurrent resolutions miss the cache for the same agent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should collapse them into a single registry call", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)

			var storeMutex sync.Mutex
			stored := ""
			release := make(chan struct{})
			misses := make(chan struct{}, 3)

			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any()).
				DoAndReturn(func(string) (string, error) {
					storeMutex.Lock()
					defer storeMutex.Unlock()
					if stored == "" {
						misses <- struct{}{}
						return "", errors.New("miss")
					}
					return stored, nil
				}).
				AnyTimes()
			suite.mockAgentRegistryRequester.EXPECT().
				GetAgent(gomock.Any(), gomock.Any()).
				DoAndReturn(func(string, uint64) (*models.RegistryAgent, error) {
					<-release
					return registryAgentForTest(), nil
				}).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(key, value string, ttl *time.Duration) error {
					storeMutex.Lock()
					defer storeMutex.Unlock()
					stored = value
					return nil
				}).
				Times(2)

			var wg sync.WaitGroup
			results := make([]*models.RegistryAgent, 2)
			for i := range results {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					agent, err := suite.sut.ResolveByCollection("0xabc", 1)
					assert.NoError(t, err)
					results[i] = agent
				}(i)
			}

			// Hold the fetch until both outer reads plus the lock holder's re-check
			// have missed; the waiter's re-check then hits the cache.
			<-misses
			<-misses
			<-misses
			close(release)
			wg.Wait()

			assert.Equal(t, results[0].Name, results[1].Name)
		})
	})

	// Writing both keys means the second lookup path costs no registry call.
	t.Run("Given the agent was resolved by ENS and is then looked up by collection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should serve the collection lookup from the entry the ENS resolve warmed", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRegistryAgentBeforeEach(t)

			var mutex sync.Mutex
			store := map[string]string{}
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any()).
				DoAndReturn(func(key string) (string, error) {
					mutex.Lock()
					defer mutex.Unlock()
					if value, ok := store[key]; ok {
						return value, nil
					}
					return "", errors.New("miss")
				}).
				AnyTimes()
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(key, value string, ttl *time.Duration) error {
					mutex.Lock()
					defer mutex.Unlock()
					store[key] = value
					return nil
				}).
				AnyTimes()

			suite.mockAgentRegistryRequester.EXPECT().
				GetMintedAgentByEnsDomain("hermes.agents6022.eth").
				Return(registryAgentForTest(), nil).
				Times(1)
			suite.mockAgentRegistryRequester.EXPECT().GetAgent(gomock.Any(), gomock.Any()).Times(0)

			byEns, err := suite.sut.ResolveByEnsDomain("hermes.agents6022.eth")
			assert.NoError(t, err)

			byCollection, err := suite.sut.ResolveByCollection("0xABC", 1)
			assert.NoError(t, err)

			assert.Equal(t, byEns.Name, byCollection.Name)
			assert.Contains(t, store, registryAgentCollectionCachePath)
		})
	})
}
