package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenStartingBlockEventsDequeueDaemonTestingSuite struct {
	sut *services.BlockEventsDequeueDaemon

	mockBlockEventsQueue     *interfaces_mocks.MockBlockEventsQueueInterface
	mockBlockEventsProcessor *interfaces_mocks.MockBlockEventsProcessorInterface
}

func WhenStartingBlockEventsDequeueDaemonBeforeEach(t *testing.T) *WhenStartingBlockEventsDequeueDaemonTestingSuite {
	mockController := gomock.NewController(t)

	mockBlockEventsQueue := interfaces_mocks.NewMockBlockEventsQueueInterface(mockController)
	mockBlockEventsProcessor := interfaces_mocks.NewMockBlockEventsProcessorInterface(mockController)

	return &WhenStartingBlockEventsDequeueDaemonTestingSuite{
		sut: services.NewBlockEventsDequeueDaemon(
			zap.NewNop(),
			&settings.SyncSettings{BlockEventsPollInterval: 5 * time.Millisecond},
			mockBlockEventsQueue,
			mockBlockEventsProcessor,
		),
		mockBlockEventsQueue:     mockBlockEventsQueue,
		mockBlockEventsProcessor: mockBlockEventsProcessor,
	}
}

func TestWhenStartingBlockEventsDequeueDaemon(t *testing.T) {
	t.Parallel()

	t.Run("Given queued block events", func(t *testing.T) {
		t.Parallel()

		t.Run("Should keep processing them without blocking", func(t *testing.T) {
			t.Parallel()

			suite := WhenStartingBlockEventsDequeueDaemonBeforeEach(t)
			item := &interfaces.QueueItem{BlockNumber: 100}
			processed := make(chan struct{}, 100)
			suite.mockBlockEventsQueue.EXPECT().Dequeue(gomock.Any()).Return(uint64(100), item).MinTimes(2)
			suite.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), uint64(100), item).
				DoAndReturn(func(context.Context, uint64, *interfaces.QueueItem) error {
					processed <- struct{}{}
					return nil
				}).
				MinTimes(2)

			suite.sut.Start(context.Background())
			<-processed
			<-processed
			suite.sut.Stop()

			assert.Empty(t, suite.sut.Failed())
		})
	})

	t.Run("Given processing fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report the error and process nothing more", func(t *testing.T) {
			t.Parallel()

			suite := WhenStartingBlockEventsDequeueDaemonBeforeEach(t)
			item := &interfaces.QueueItem{BlockNumber: 100}
			suite.mockBlockEventsQueue.EXPECT().Dequeue(gomock.Any()).Return(uint64(100), item).Times(1)
			suite.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), uint64(100), item).
				Return(assert.AnError).
				Times(1)

			suite.sut.Start(context.Background())

			assert.Equal(t, assert.AnError, <-suite.sut.Failed())
			time.Sleep(20 * time.Millisecond)
			suite.sut.Stop()
		})
	})

	t.Run("Given processing is cancelled while the daemon runs", func(t *testing.T) {
		t.Parallel()

		t.Run("Should report it", func(t *testing.T) {
			t.Parallel()

			suite := WhenStartingBlockEventsDequeueDaemonBeforeEach(t)
			item := &interfaces.QueueItem{BlockNumber: 100}
			suite.mockBlockEventsQueue.EXPECT().Dequeue(gomock.Any()).Return(uint64(100), item).Times(1)
			suite.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), uint64(100), item).
				Return(context.Canceled).
				Times(1)

			suite.sut.Start(context.Background())

			assert.ErrorIs(t, <-suite.sut.Failed(), context.Canceled)
			suite.sut.Stop()
		})
	})

	t.Run("Given a daemon started twice", func(t *testing.T) {
		t.Parallel()

		t.Run("Should run a single loop", func(t *testing.T) {
			t.Parallel()

			suite := WhenStartingBlockEventsDequeueDaemonBeforeEach(t)
			item := &interfaces.QueueItem{BlockNumber: 100}
			suite.mockBlockEventsQueue.EXPECT().Dequeue(gomock.Any()).Return(uint64(100), item).Times(1)
			suite.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), uint64(100), item).
				DoAndReturn(func(ctx context.Context, _ uint64, _ *interfaces.QueueItem) error {
					<-ctx.Done()
					return ctx.Err()
				}).
				Times(1)

			suite.sut.Start(context.Background())
			suite.sut.Start(context.Background())
			time.Sleep(20 * time.Millisecond)
			suite.sut.Stop()

			assert.Empty(t, suite.sut.Failed())
		})
	})
}
