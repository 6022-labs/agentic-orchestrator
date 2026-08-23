package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenDequeuingBlockEventsTestingSuite struct {
	sut *services.BlockEventsDequeueDaemon

	mockBlockEventsQueue     *interfaces_mocks.MockBlockEventsQueueInterface
	mockBlockEventsProcessor *interfaces_mocks.MockBlockEventsProcessorInterface
}

func WhenDequeuingBlockEventsBeforeEach(t *testing.T) *WhenDequeuingBlockEventsTestingSuite {
	mockController := gomock.NewController(t)

	mockBlockEventsQueue := interfaces_mocks.NewMockBlockEventsQueueInterface(mockController)
	mockBlockEventsProcessor := interfaces_mocks.NewMockBlockEventsProcessorInterface(mockController)

	sut := services.NewBlockEventsDequeueDaemon(
		zap.NewNop(),
		mockBlockEventsQueue,
		mockBlockEventsProcessor,
	)

	return &WhenDequeuingBlockEventsTestingSuite{
		sut:                      sut,
		mockBlockEventsQueue:     mockBlockEventsQueue,
		mockBlockEventsProcessor: mockBlockEventsProcessor,
	}
}

func TestWhenDequeuingBlockEvents(t *testing.T) {
	t.Parallel()

	t.Run("Given the queue is empty", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenDequeuingBlockEventsTestingSuite) {
			suite.mockBlockEventsQueue.EXPECT().
				Dequeue(gomock.Any()).
				Return(uint64(0), nil).
				Times(1)
		}

		t.Run("Should return nil without processing", func(t *testing.T) {
			t.Parallel()
			basics := WhenDequeuingBlockEventsBeforeEach(t)
			initSuite(basics)

			err := basics.sut.DequeueAndProcess(context.Background())

			assert.Nil(t, err)
		})
	})

	t.Run("Given the queue has an item and processing succeeds", func(t *testing.T) {
		t.Parallel()

		t.Run("Should process the block events successfully", func(t *testing.T) {
			t.Parallel()
			basics := WhenDequeuingBlockEventsBeforeEach(t)

			blockNumber := uint64(100)
			queueItem := &interfaces.QueueItem{
				BlockNumber: blockNumber,
			}

			basics.mockBlockEventsQueue.EXPECT().
				Dequeue(gomock.Any()).
				Return(blockNumber, queueItem).
				Times(1)

			basics.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), blockNumber, queueItem).
				Return(nil).
				Times(1)

			err := basics.sut.DequeueAndProcess(context.Background())

			assert.Nil(t, err)
		})
	})

	t.Run("Given the queue has an item but processing fails", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenDequeuingBlockEventsTestingSuite) (*interfaces.QueueItem, uint64) {
			blockNumber := uint64(200)
			queueItem := &interfaces.QueueItem{
				BlockNumber: blockNumber,
			}

			suite.mockBlockEventsQueue.EXPECT().
				Dequeue(gomock.Any()).
				Return(blockNumber, queueItem).
				Times(1)

			suite.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), blockNumber, queueItem).
				Return(assert.AnError).
				Times(1)

			return queueItem, blockNumber
		}

		t.Run("Should return the processing error", func(t *testing.T) {
			t.Parallel()
			basics := WhenDequeuingBlockEventsBeforeEach(t)
			initSuite(basics)

			err := basics.sut.DequeueAndProcess(context.Background())

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given the context is cancelled before processing", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return context cancelled error", func(t *testing.T) {
			t.Parallel()
			basics := WhenDequeuingBlockEventsBeforeEach(t)

			blockNumber := uint64(300)
			queueItem := &interfaces.QueueItem{
				BlockNumber: blockNumber,
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			basics.mockBlockEventsQueue.EXPECT().
				Dequeue(gomock.Any()).
				Return(blockNumber, queueItem).
				Times(1)

			err := basics.sut.DequeueAndProcess(ctx)

			assert.NotNil(t, err)
			assert.ErrorIs(t, err, context.Canceled)
		})
	})
}
