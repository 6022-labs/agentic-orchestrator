package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	sync_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/metrics_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenQueuingBlockEventsTestingSuite struct {
	// depths receives every published (ready, pending) snapshot, so a Given can assert them.
	depths chan [2]int64

	sut *services.BlockEventsQueue
}

func WhenQueuingBlockEventsBeforeEach(t *testing.T) *WhenQueuingBlockEventsTestingSuite {
	mockSyncMetrics := sync_metrics_mocks.NewMockSyncMetricsInterface(gomock.NewController(t))

	// Captured here, not per Given: a blanket AnyTimes matches first and would swallow it.
	depths := make(chan [2]int64, 32)
	mockSyncMetrics.EXPECT().
		RecordQueueState(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, ready int64, pending int64) {
			depths <- [2]int64{ready, pending}
		}).
		AnyTimes()
	mockSyncMetrics.EXPECT().RecordHighestBlockWithEvents(gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewBlockEventsQueue(
		zap.NewNop(),
		mockSyncMetrics,
		&settings.SyncSettings{
			StartBlock:             1,
			BlockFinalizationDelay: 10 * time.Millisecond,
		},
	)

	return &WhenQueuingBlockEventsTestingSuite{
		depths: depths,

		sut: sut,
	}
}

// nextDepths fails fast instead of blocking forever when a mutation publishes nothing.
func nextDepths(t *testing.T, depths chan [2]int64, what string) [2]int64 {
	t.Helper()

	select {
	case snapshot := <-depths:
		return snapshot
	case <-time.After(2 * time.Second):
		t.Fatalf("no queue depth was published after %s", what)
		return [2]int64{}
	}
}

func waitForFinalization() {
	time.Sleep(50 * time.Millisecond)
}

func TestWhenQueuingBlockEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given the queue is empty", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil when dequeuing", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			blockNumber, item := basics.sut.Dequeue(context.Background())

			assert.Equal(t, uint64(0), blockNumber)
			assert.Nil(t, item)
		})
	})

	t.Run("Given minted events are added and block finalization timer fires", func(t *testing.T) {
		t.Parallel()

		t.Run("Should dequeue the block with events", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			mintedEvent := &abi.AgentCollectionV1Minted{}
			err := basics.sut.AddAgentCollectionV1MintedEvents(
				context.Background(),
				10,
				[]*abi.AgentCollectionV1Minted{mintedEvent},
			)
			assert.Nil(t, err)

			waitForFinalization()

			blockNumber, item := basics.sut.Dequeue(context.Background())

			assert.Equal(t, uint64(10), blockNumber)
			assert.NotNil(t, item)
			assert.Len(t, item.AgentCollectionV1MintedEvents, 1)
		})
	})

	t.Run("Given events are added but finalization timer has not fired yet", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not dequeue the block", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			mintedEvent := &abi.AgentCollectionV1Minted{}
			err := basics.sut.AddAgentCollectionV1MintedEvents(
				context.Background(),
				10,
				[]*abi.AgentCollectionV1Minted{mintedEvent},
			)
			assert.Nil(t, err)

			blockNumber, item := basics.sut.Dequeue(context.Background())

			assert.Equal(t, uint64(0), blockNumber)
			assert.Nil(t, item)
		})
	})

	t.Run("Given multiple block events are added in order", func(t *testing.T) {
		t.Parallel()

		t.Run("Should dequeue blocks in ascending order", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			err := basics.sut.AddAgentCollectionV1MintedEvents(
				context.Background(),
				20,
				[]*abi.AgentCollectionV1Minted{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1TransferEvents(
				context.Background(),
				10,
				[]*abi.AgentCollectionV1Transfer{{}},
			)
			assert.Nil(t, err)

			waitForFinalization()

			blockNumber1, item1 := basics.sut.Dequeue(context.Background())
			assert.Equal(t, uint64(10), blockNumber1)
			assert.NotNil(t, item1)
			assert.Len(t, item1.AgentCollectionV1TransferEvents, 1)

			blockNumber2, item2 := basics.sut.Dequeue(context.Background())
			assert.Equal(t, uint64(20), blockNumber2)
			assert.NotNil(t, item2)
			assert.Len(t, item2.AgentCollectionV1MintedEvents, 1)
		})
	})

	t.Run("Given multiple event types are added to the same block", func(t *testing.T) {
		t.Parallel()

		t.Run("Should aggregate all events in a single queue item", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			blockNumber := uint64(50)

			err := basics.sut.AddAgentCollectionV1MintedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1Minted{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1TransferEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1Transfer{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1AgentImageUpdatedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1AgentImageUpdated{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1AgentImageRemovedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1AgentImageRemoved{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1MintProposalCreatedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1MintProposalCreated{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1MintProposalRemovedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1MintProposalRemoved{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1AgentAttributeUpdatedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionV1AgentAttributeUpdated{{}},
			)
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionsManagerCollectionCreatedEvents(
				context.Background(),
				blockNumber,
				[]*abi.AgentCollectionsManagerCollectionCreated{{}},
			)
			assert.Nil(t, err)

			waitForFinalization()

			dequeuedBlockNumber, item := basics.sut.Dequeue(context.Background())

			assert.Equal(t, blockNumber, dequeuedBlockNumber)
			assert.NotNil(t, item)
			assert.Len(t, item.AgentCollectionV1MintedEvents, 1)
			assert.Len(t, item.AgentCollectionV1TransferEvents, 1)
			assert.Len(t, item.AgentCollectionV1AgentImageUpdatedEvents, 1)
			assert.Len(t, item.AgentCollectionV1AgentImageRemovedEvents, 1)
			assert.Len(t, item.AgentCollectionV1MintProposalCreatedEvents, 1)
			assert.Len(t, item.AgentCollectionV1MintProposalRemovedEvents, 1)
			assert.Len(t, item.AgentCollectionV1AgentAttributeUpdatedEvents, 1)
			assert.Len(t, item.AgentCollectionsManagerCollectionCreatedEvents, 1)
		})
	})

	t.Run("Given a block is dequeued", func(t *testing.T) {
		t.Parallel()

		t.Run("Should remove it from the queue", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			err := basics.sut.AddAgentCollectionV1MintedEvents(
				context.Background(),
				10,
				[]*abi.AgentCollectionV1Minted{{}},
			)
			assert.Nil(t, err)

			waitForFinalization()

			blockNumber, item := basics.sut.Dequeue(context.Background())
			assert.Equal(t, uint64(10), blockNumber)
			assert.NotNil(t, item)

			blockNumber2, item2 := basics.sut.Dequeue(context.Background())
			assert.Equal(t, uint64(0), blockNumber2)
			assert.Nil(t, item2)
		})
	})

	t.Run("Given the published queue depths are observed", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report a block as pending, then ready, then drained", func(t *testing.T) {
			t.Parallel()

			basics := WhenQueuingBlockEventsBeforeEach(t)

			err := basics.sut.AddAgentCollectionV1MintedEvents(
				context.Background(),
				2,
				[]*abi.AgentCollectionV1Minted{{}},
			)
			assert.Nil(t, err)
			assert.Equal(
				t,
				[2]int64{0, 1},
				nextDepths(t, basics.depths, "an add"),
				"a freshly added block is pending, not ready",
			)

			// The finalization timer moves it to ready and must publish that itself.
			waitForFinalization()
			assert.Equal(t, [2]int64{1, 0}, nextDepths(t, basics.depths, "the finalization timer"),
				"the finalization timer publishes the new depths")

			blockNumber, item := basics.sut.Dequeue(context.Background())
			assert.Equal(t, uint64(2), blockNumber)
			assert.NotNil(t, item)
			assert.Equal(t, [2]int64{0, 0}, nextDepths(t, basics.depths, "a dequeue"), "dequeuing drains the queue")
		})
	})
}
