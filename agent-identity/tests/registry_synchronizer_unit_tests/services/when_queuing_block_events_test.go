package services_test

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type WhenQueuingBlockEventsTestingSuite struct {
	sut *services.BlockEventsQueue
}

func WhenQueuingBlockEventsBeforeEach(t *testing.T) *WhenQueuingBlockEventsTestingSuite {
	sut := services.NewBlockEventsQueue(
		zap.NewNop(),
		&settings.SyncSettings{
			StartBlock:             1,
			BlockFinalizationDelay: 10 * time.Millisecond,
		},
	)

	return &WhenQueuingBlockEventsTestingSuite{
		sut: sut,
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

			blockNumber, item := basics.sut.Dequeue()

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
			err := basics.sut.AddAgentCollectionV1MintedEvents(10, []*abi.AgentCollectionV1Minted{mintedEvent})
			assert.Nil(t, err)

			waitForFinalization()

			blockNumber, item := basics.sut.Dequeue()

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
			err := basics.sut.AddAgentCollectionV1MintedEvents(10, []*abi.AgentCollectionV1Minted{mintedEvent})
			assert.Nil(t, err)

			blockNumber, item := basics.sut.Dequeue()

			assert.Equal(t, uint64(0), blockNumber)
			assert.Nil(t, item)
		})
	})

	t.Run("Given multiple block events are added in order", func(t *testing.T) {
		t.Parallel()

		t.Run("Should dequeue blocks in ascending order", func(t *testing.T) {
			t.Parallel()
			basics := WhenQueuingBlockEventsBeforeEach(t)

			err := basics.sut.AddAgentCollectionV1MintedEvents(20, []*abi.AgentCollectionV1Minted{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1TransferEvents(10, []*abi.AgentCollectionV1Transfer{{}})
			assert.Nil(t, err)

			waitForFinalization()

			blockNumber1, item1 := basics.sut.Dequeue()
			assert.Equal(t, uint64(10), blockNumber1)
			assert.NotNil(t, item1)
			assert.Len(t, item1.AgentCollectionV1TransferEvents, 1)

			blockNumber2, item2 := basics.sut.Dequeue()
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

			err := basics.sut.AddAgentCollectionV1MintedEvents(blockNumber, []*abi.AgentCollectionV1Minted{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1TransferEvents(blockNumber, []*abi.AgentCollectionV1Transfer{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1AgentImageUpdatedEvents(blockNumber, []*abi.AgentCollectionV1AgentImageUpdated{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1MintProposalCreatedEvents(blockNumber, []*abi.AgentCollectionV1MintProposalCreated{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1MintProposalRemovedEvents(blockNumber, []*abi.AgentCollectionV1MintProposalRemoved{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionV1AgentAttributeUpdatedEvents(blockNumber, []*abi.AgentCollectionV1AgentAttributeUpdated{{}})
			assert.Nil(t, err)

			err = basics.sut.AddAgentCollectionsManagerCollectionCreatedEvents(blockNumber, []*abi.AgentCollectionsManagerCollectionCreated{{}})
			assert.Nil(t, err)

			waitForFinalization()

			dequeuedBlockNumber, item := basics.sut.Dequeue()

			assert.Equal(t, blockNumber, dequeuedBlockNumber)
			assert.NotNil(t, item)
			assert.Len(t, item.AgentCollectionV1MintedEvents, 1)
			assert.Len(t, item.AgentCollectionV1TransferEvents, 1)
			assert.Len(t, item.AgentCollectionV1AgentImageUpdatedEvents, 1)
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

			err := basics.sut.AddAgentCollectionV1MintedEvents(10, []*abi.AgentCollectionV1Minted{{}})
			assert.Nil(t, err)

			waitForFinalization()

			blockNumber, item := basics.sut.Dequeue()
			assert.Equal(t, uint64(10), blockNumber)
			assert.NotNil(t, item)

			blockNumber2, item2 := basics.sut.Dequeue()
			assert.Equal(t, uint64(0), blockNumber2)
			assert.Nil(t, item2)
		})
	})
}
