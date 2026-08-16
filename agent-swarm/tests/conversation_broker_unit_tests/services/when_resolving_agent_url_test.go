package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenResolvingAgentUrlTestingSuite struct {
	sut *services.AgentUrlResolver

	mockCacheStorage   *interfaces_mocks.MockCacheStorageInterface
	mockEnsUrlResolver *interfaces_mocks.MockEnsUrlResolverInterface
}

func WhenResolvingAgentUrlBeforeEach(t *testing.T) *WhenResolvingAgentUrlTestingSuite {
	mockController := gomock.NewController(t)

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockEnsUrlResolver := interfaces_mocks.NewMockEnsUrlResolverInterface(mockController)

	sut := services.NewAgentUrlResolver(
		zap.NewNop(),
		mockCacheStorage,
		mockEnsUrlResolver,
		&settings.RemoteAgentResolutionSettings{CacheTtlSeconds: 300},
	)

	return &WhenResolvingAgentUrlTestingSuite{
		sut: sut,

		mockCacheStorage:   mockCacheStorage,
		mockEnsUrlResolver: mockEnsUrlResolver,
	}
}

func TestWhenResolvingAgentUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given the url is not cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read("agent-urls/agent.eth").Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("https://node.example", nil).Times(1)
			suite.mockCacheStorage.EXPECT().
				Store("agent-urls/agent.eth", "https://node.example", gomock.Any()).
				DoAndReturn(func(key, value string, ttl *time.Duration) error {
					assert.Equal(t, 300*time.Second, *ttl)
					return nil
				}).
				Times(1)
		}

		t.Run("Should resolve the url and cache it with the TTL", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			url, err := suite.sut.ResolveUrl("Agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "https://node.example", url)
		})
	})

	t.Run("Given the url is cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read("agent-urls/agent.eth").Return("https://node.example", nil).Times(1)
		}

		t.Run("Should serve the url from the cache without resolving", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			url, err := suite.sut.ResolveUrl("agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "https://node.example", url)
		})
	})

	t.Run("Given the ENS resolution fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("", errors.New("boom")).Times(1)
		}

		t.Run("Should propagate the error and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
		})
	})

	t.Run("Given no url is published", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("", nil).Times(1)
		}

		t.Run("Should return an error and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
		})
	})

	t.Run("Given caching the url fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingAgentUrlTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockEnsUrlResolver.EXPECT().ResolveUrl("agent.eth").Return("https://node.example", nil).Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(errors.New("boom")).
				Times(1)
		}

		t.Run("Should still return the resolved url", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingAgentUrlBeforeEach(t)
			initSuite(suite)

			url, err := suite.sut.ResolveUrl("agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "https://node.example", url)
		})
	})
}
