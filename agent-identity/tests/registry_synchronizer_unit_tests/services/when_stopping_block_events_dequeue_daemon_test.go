package services_test

import (
	"context"
	"sync/atomic"
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

type WhenStoppingBlockEventsDequeueDaemonTestingSuite struct {
	sut *services.BlockEventsDequeueDaemon

	mockBlockEventsQueue     *interfaces_mocks.MockBlockEventsQueueInterface
	mockBlockEventsProcessor *interfaces_mocks.MockBlockEventsProcessorInterface
}

func WhenStoppingBlockEventsDequeueDaemonBeforeEach(t *testing.T) *WhenStoppingBlockEventsDequeueDaemonTestingSuite {
	mockController := gomock.NewController(t)

	mockBlockEventsQueue := interfaces_mocks.NewMockBlockEventsQueueInterface(mockController)
	mockBlockEventsProcessor := interfaces_mocks.NewMockBlockEventsProcessorInterface(mockController)

	return &WhenStoppingBlockEventsDequeueDaemonTestingSuite{
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

func TestWhenStoppingBlockEventsDequeueDaemon(t *testing.T) {
	t.Parallel()

	t.Run("Given a daemon never started", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return at once", func(t *testing.T) {
			t.Parallel()

			assert.NotPanics(t, WhenStoppingBlockEventsDequeueDaemonBeforeEach(t).sut.Stop)
		})
	})

	t.Run("Given a block in progress", func(t *testing.T) {
		t.Parallel()

		t.Run("Should wait for it and report no failure", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoppingBlockEventsDequeueDaemonBeforeEach(t)
			item := &interfaces.QueueItem{BlockNumber: 100}
			entered := make(chan struct{})
			finished := &atomic.Bool{}
			suite.mockBlockEventsQueue.EXPECT().Dequeue(gomock.Any()).Return(uint64(100), item).Times(1)
			suite.mockBlockEventsProcessor.EXPECT().
				ProcessBlock(gomock.Any(), uint64(100), item).
				DoAndReturn(func(ctx context.Context, _ uint64, _ *interfaces.QueueItem) error {
					close(entered)
					<-ctx.Done()
					time.Sleep(20 * time.Millisecond)
					finished.Store(true)
					return ctx.Err()
				}).
				Times(1)

			suite.sut.Start(context.Background())
			<-entered
			suite.sut.Stop()

			assert.True(t, finished.Load())
			assert.Empty(t, suite.sut.Failed())
		})
	})
}
