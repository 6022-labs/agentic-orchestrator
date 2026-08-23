package services_test

import (
	"context"

	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/traces/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/model_builders"
	broker_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/services_mocks/interfaces_mocks"
	broker_traces_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_mocks/traces_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenResolvingRemoteAgentCardTestingSuite struct {
	sut *services.RemoteAgentCardResolver

	mockCacheStorage               *interfaces_mocks.MockCacheStorageInterface
	mockAgentUrlResolver           *interfaces_mocks.MockAgentUrlResolverInterface
	mockRegistryAgentResolver      *interfaces_mocks.MockRegistryAgentResolverInterface
	mockRemoteAgentCardRequester   *interfaces_mocks.MockRemoteAgentCardRequesterInterface
	mockAgentCardSignatureVerifier *interfaces_mocks.MockAgentCardSignatureVerifierInterface
}

func WhenResolvingRemoteAgentCardBeforeEach(t *testing.T) *WhenResolvingRemoteAgentCardTestingSuite {
	mockController := gomock.NewController(t)

	mockBrokerMetrics := broker_metrics_mocks.NewMockBrokerMetricsInterface(mockController)
	mockBrokerMetrics.EXPECT().RecordCacheResult(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockBrokerMetrics.EXPECT().RecordAgentResolution(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	mockBrokerTracer := broker_traces_mocks.NewMockBrokerTracerInterface(mockController)
	mockSpan := broker_traces_mocks.NewMockSpan(mockController)
	mockSpan.EXPECT().End().AnyTimes()
	mockSpan.EXPECT().Fail(gomock.Any()).AnyTimes()
	mockBrokerTracer.EXPECT().
		StartAgentResolution(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ string, _ string) (context.Context, traces_interfaces.Span) {
			return ctx, mockSpan
		}).
		AnyTimes()

	mockCacheStorage := interfaces_mocks.NewMockCacheStorageInterface(mockController)
	mockAgentUrlResolver := interfaces_mocks.NewMockAgentUrlResolverInterface(mockController)
	mockRegistryAgentResolver := interfaces_mocks.NewMockRegistryAgentResolverInterface(mockController)
	mockRemoteAgentCardRequester := interfaces_mocks.NewMockRemoteAgentCardRequesterInterface(mockController)
	mockAgentCardSignatureVerifier := interfaces_mocks.NewMockAgentCardSignatureVerifierInterface(mockController)

	sut := services.NewRemoteAgentCardResolver(
		zap.NewNop(),
		mockBrokerMetrics,
		mockBrokerTracer,
		mockCacheStorage,
		mockAgentUrlResolver,
		mockRegistryAgentResolver,
		mockRemoteAgentCardRequester,
		mockAgentCardSignatureVerifier,
		&settings.RemoteAgentResolutionSettings{CacheTtlSeconds: 300},
	)

	return &WhenResolvingRemoteAgentCardTestingSuite{
		sut: sut,

		mockCacheStorage:               mockCacheStorage,
		mockAgentUrlResolver:           mockAgentUrlResolver,
		mockRegistryAgentResolver:      mockRegistryAgentResolver,
		mockRemoteAgentCardRequester:   mockRemoteAgentCardRequester,
		mockAgentCardSignatureVerifier: mockAgentCardSignatureVerifier,
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

// rawCardBody is what the node served; verification runs over these bytes, not the parsed struct.
const rawCardBody = `{"protocol":"6022","endpoints":{"responses":"https://node.example/api/responses"}}`

// signedCard is a node that made an integrity claim; cardWithResponses is one that made none.
func signedCard() *models.RemoteAgentCard {
	card := cardWithResponses()
	card.Signatures = []models.AgentCardSignature{{Protected: "protected", Signature: "signature"}}
	return card
}

func TestWhenResolvingRemoteAgentCard(t *testing.T) {
	t.Parallel()

	t.Run("Given the card is not cached", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return("", errors.New("miss")).
				Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("https://node.example", nil).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(cardWithResponses(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), "remote-agent-cards/agent.eth", gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, key, value string, ttl *time.Duration) error {
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

			card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

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
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				DoAndReturn(func(context.Context, string) (string, error) {
					storeMutex.Lock()
					defer storeMutex.Unlock()
					if stored == "" {
						misses <- struct{}{}
						return "", errors.New("miss")
					}
					return stored, nil
				}).
				Times(4)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				DoAndReturn(func(context.Context, string) (string, error) {
					<-release
					return "https://node.example", nil
				}).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(cardWithResponses(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, key, value string, ttl *time.Duration) error {
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
					card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")
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
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return(string(cached), nil).
				Times(1)
		}

		t.Run("Should serve the card from the cache without resolving", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
		})
	})

	t.Run("Given the cached value is corrupted", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return("not-json", nil).
				Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("https://node.example", nil).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(cardWithResponses(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
		}

		t.Run("Should fall back to a full resolution", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
		})
	})

	t.Run("Given the card carries no signature", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return("", errors.New("miss")).
				Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("https://node.example", nil).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(cardWithResponses(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
		}

		t.Run("Should accept it without resolving a signer at all", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)
			// Integrity is the node's responsibility; making no claim costs no registry read.
			suite.mockRegistryAgentResolver.EXPECT().ResolveByEnsDomain(gomock.Any(), gomock.Any()).Times(0)
			suite.mockAgentCardSignatureVerifier.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
		})
	})

	t.Run("Given the card carries a signature that verifies", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return("", errors.New("miss")).
				Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("https://node.example", nil).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(signedCard(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any(), "agent.eth").
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithAddress(models.EvmAddressType, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").
						Build(),
					nil,
				).
				Times(1)
			suite.mockAgentCardSignatureVerifier.EXPECT().
				Verify([]byte(rawCardBody), signedCard().Signatures, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").
				Return(nil).
				Times(1)
			suite.mockCacheStorage.EXPECT().
				Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1)
		}

		t.Run("Should verify against the registry signer and cache the card", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.NoError(t, err)
			assert.Equal(t, "6022", card.Protocol)
		})
	})

	t.Run("Given the card carries a signature that does not verify", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return("", errors.New("miss")).
				Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("https://node.example", nil).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(signedCard(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any(), "agent.eth").
				Return(
					model_builders.NewRegistryAgentBuilder().
						WithAddress(models.EvmAddressType, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266").
						Build(),
					nil,
				).
				Times(1)
			suite.mockAgentCardSignatureVerifier.EXPECT().
				Verify(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(errors.New("signer mismatch")).
				Times(1)
		}

		t.Run("Should fail and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)
			suite.mockCacheStorage.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			card, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.Error(t, err)
			assert.Nil(t, card)
		})
	})

	t.Run("Given a signed card whose signer cannot be resolved", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().
				Read(gomock.Any(), "remote-agent-cards/agent.eth").
				Return("", errors.New("miss")).
				Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("https://node.example", nil).
				Times(1)
			suite.mockRemoteAgentCardRequester.EXPECT().
				GetCard(gomock.Any(), "https://node.example").
				Return(signedCard(), []byte(rawCardBody), nil).
				Times(1)
			suite.mockRegistryAgentResolver.EXPECT().
				ResolveByEnsDomain(gomock.Any(), "agent.eth").
				Return(nil, errors.New("registry down")).
				Times(1)
		}

		t.Run("Should fail", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)
			suite.mockAgentCardSignatureVerifier.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			_, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.Error(t, err)
		})
	})

	t.Run("Given the ENS resolution fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenResolvingRemoteAgentCardTestingSuite) {
			suite.mockCacheStorage.EXPECT().Read(gomock.Any(), gomock.Any()).Return("", errors.New("miss")).Times(2)
			suite.mockAgentUrlResolver.EXPECT().
				ResolveUrl(gomock.Any(), "agent.eth").
				Return("", errors.New("boom")).
				Times(1)
		}

		t.Run("Should propagate the error and cache nothing", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingRemoteAgentCardBeforeEach(t)
			initSuite(suite)

			_, err := suite.sut.ResolveCard(context.Background(), "agent.eth")

			assert.Error(t, err)
		})
	})
}
