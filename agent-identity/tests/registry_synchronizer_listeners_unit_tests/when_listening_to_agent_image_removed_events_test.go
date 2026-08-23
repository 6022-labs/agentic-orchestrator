package registry_synchronizer_listeners_unit_tests

import (
	"context"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_listeners"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/repositories_mocks"
	chain_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_listeners_mocks/metrics_mocks/interfaces_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenListeningToAgentImageRemovedEventsTestingSuite struct {
	sut *registry_synchronizer_listeners.AgentCollectionAgentImageRemovedListener

	chainEventMetricsMock *chain_metrics_mocks.MockChainEventMetricsInterface
	subscriptionProvider  *interfaces_mocks.MockAgentCollectionEventSubscriptionProviderInterface
	blockEventsQueue      *interfaces_mocks.MockBlockEventsQueueInterface

	events       chan *abi.AgentCollectionV1AgentImageRemoved
	subscription *stubSubscription
}

func WhenListeningToAgentImageRemovedEventsBeforeEach(
	t *testing.T,
) *WhenListeningToAgentImageRemovedEventsTestingSuite {
	mockController := gomock.NewController(t)

	chainEventMetricsMock := chain_metrics_mocks.NewMockChainEventMetricsInterface(mockController)
	agentCollectionRepository := repositories_mocks.NewMockAgentCollectionRepositoryInterface(mockController)
	blockEventsQueue := interfaces_mocks.NewMockBlockEventsQueueInterface(mockController)
	subscriptionProvider := interfaces_mocks.NewMockAgentCollectionEventSubscriptionProviderInterface(mockController)

	events := make(chan *abi.AgentCollectionV1AgentImageRemoved, 1)
	subscription := newStubSubscription()

	agentCollectionRepository.EXPECT().
		GetAll(gomock.Any()).
		Return([]*models.AgentCollection{{Address: "0x1111111111111111111111111111111111111111"}}, nil)

	subscriptionProvider.EXPECT().
		StartAgentImageRemovedSubscription(gomock.Any(), gomock.Any()).
		Return((<-chan *abi.AgentCollectionV1AgentImageRemoved)(events), subscription, nil)

	sut := registry_synchronizer_listeners.NewAgentCollectionAgentImageRemovedListener(
		zap.NewNop(),
		chainEventMetricsMock,
		&common_settings.BlockchainSettings{ChainId: int64(testChainId)},
		agentCollectionRepository,
		blockEventsQueue,
		subscriptionProvider,
	)

	return &WhenListeningToAgentImageRemovedEventsTestingSuite{
		sut: sut,

		chainEventMetricsMock: chainEventMetricsMock,
		subscriptionProvider:  subscriptionProvider,
		blockEventsQueue:      blockEventsQueue,

		events:       events,
		subscription: subscription,
	}
}

func newAgentImageRemovedEvent(blockNumber uint64) *abi.AgentCollectionV1AgentImageRemoved {
	return &abi.AgentCollectionV1AgentImageRemoved{Raw: types.Log{BlockNumber: blockNumber}}
}

func TestWhenListeningToAgentImageRemovedEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given the subscription is opened and then closed", func(t *testing.T) {
		t.Parallel()

		t.Run("Should leave no subscription counted as open", func(t *testing.T) {
			t.Parallel()

			suite := WhenListeningToAgentImageRemovedEventsBeforeEach(t)

			opened := make(chan struct{})
			closed := make(chan struct{})
			suite.chainEventMetricsMock.EXPECT().
				RecordSubscriptionOpened(gomock.Any(), "AgentImageRemoved", testChainId).
				Do(func(context.Context, string, uint64) { close(opened) })
			suite.chainEventMetricsMock.EXPECT().
				RecordSubscriptionClosed(gomock.Any(), "AgentImageRemoved", testChainId).
				Do(func(context.Context, string, uint64) { close(closed) })

			ctx, cancel := context.WithCancel(context.Background())
			assert.NoError(t, suite.sut.SubscribeAll(ctx))

			select {
			case <-opened:
			case <-time.After(time.Second):
				t.Fatal("subscription was never counted as opened")
			}

			cancel()

			select {
			case <-closed:
			case <-time.After(time.Second):
				t.Fatal("subscription was never counted as closed")
			}
		})
	})

	t.Run("Given an event is received", func(t *testing.T) {
		t.Parallel()

		t.Run("Should record it as handled once enqueued", func(t *testing.T) {
			t.Parallel()

			suite := WhenListeningToAgentImageRemovedEventsBeforeEach(t)

			suite.chainEventMetricsMock.EXPECT().
				RecordSubscriptionOpened(gomock.Any(), "AgentImageRemoved", testChainId)

			closed := make(chan struct{})
			suite.chainEventMetricsMock.EXPECT().
				RecordSubscriptionClosed(gomock.Any(), "AgentImageRemoved", testChainId).
				Do(func(context.Context, string, uint64) { close(closed) })

			suite.blockEventsQueue.EXPECT().
				AddAgentCollectionV1AgentImageRemovedEvents(gomock.Any(), uint64(42), gomock.Len(1)).
				Return(nil)

			recorded := make(chan string, 1)
			suite.chainEventMetricsMock.EXPECT().
				RecordEvent(gomock.Any(), "AgentImageRemoved", testChainId, gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _ string, _ uint64, outcome string, _ time.Duration) {
					recorded <- outcome
				})

			ctx, cancel := context.WithCancel(context.Background())
			assert.NoError(t, suite.sut.SubscribeAll(ctx))

			suite.events <- newAgentImageRemovedEvent(42)

			select {
			case outcome := <-recorded:
				assert.Equal(t, "handled", outcome)
			case <-time.After(time.Second):
				t.Fatal("event was never recorded")
			}

			// Drain the watcher, or its deferred close races the controller's Finish.
			cancel()
			select {
			case <-closed:
			case <-time.After(time.Second):
				t.Fatal("subscription was never counted as closed")
			}
		})
	})

	t.Run("Given enqueuing the event fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should record it as failed", func(t *testing.T) {
			t.Parallel()

			suite := WhenListeningToAgentImageRemovedEventsBeforeEach(t)

			suite.chainEventMetricsMock.EXPECT().
				RecordSubscriptionOpened(gomock.Any(), "AgentImageRemoved", testChainId)

			closed := make(chan struct{})
			suite.chainEventMetricsMock.EXPECT().
				RecordSubscriptionClosed(gomock.Any(), "AgentImageRemoved", testChainId).
				Do(func(context.Context, string, uint64) { close(closed) })

			suite.blockEventsQueue.EXPECT().
				AddAgentCollectionV1AgentImageRemovedEvents(gomock.Any(), uint64(43), gomock.Len(1)).
				Return(assert.AnError)

			recorded := make(chan string, 1)
			suite.chainEventMetricsMock.EXPECT().
				RecordEvent(gomock.Any(), "AgentImageRemoved", testChainId, gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _ string, _ uint64, outcome string, _ time.Duration) {
					recorded <- outcome
				})

			ctx, cancel := context.WithCancel(context.Background())
			assert.NoError(t, suite.sut.SubscribeAll(ctx))

			suite.events <- newAgentImageRemovedEvent(43)

			select {
			case outcome := <-recorded:
				assert.Equal(t, "failed", outcome)
			case <-time.After(time.Second):
				t.Fatal("event was never recorded")
			}

			// Drain the watcher, or its deferred close races the controller's Finish.
			cancel()
			select {
			case <-closed:
			case <-time.After(time.Second):
				t.Fatal("subscription was never counted as closed")
			}
		})
	})
}
