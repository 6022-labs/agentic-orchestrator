package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"go.uber.org/zap"
)

type BlockEventsDequeueDaemon struct {
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce sync.Once
	stopOnce  sync.Once
	loopWg    sync.WaitGroup
	failed    chan error

	logger *zap.Logger

	syncSettings *settings.SyncSettings

	blockEventsQueue interfaces.BlockEventsQueueInterface

	blockEventsProcessor interfaces.BlockEventsProcessorInterface
}

func NewBlockEventsDequeueDaemon(
	logger *zap.Logger,
	syncSettings *settings.SyncSettings,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	blockEventsProcessor interfaces.BlockEventsProcessorInterface,
) *BlockEventsDequeueDaemon {
	return &BlockEventsDequeueDaemon{
		failed:               make(chan error, 1),
		logger:               logger,
		syncSettings:         syncSettings,
		blockEventsQueue:     blockEventsQueue,
		blockEventsProcessor: blockEventsProcessor,
	}
}

func (d *BlockEventsDequeueDaemon) Start(ctx context.Context) {
	d.startOnce.Do(func() {
		d.ctx, d.cancel = context.WithCancel(ctx)

		d.logger.Info("Starting block event dequeue synchronizer daemon...")
		d.loopWg.Add(1)
		go d.runLoop()
	})
}

func (d *BlockEventsDequeueDaemon) Stop() {
	d.stopOnce.Do(func() {
		d.logger.Info("Stopping block event dequeue synchronizer daemon...")
		if d.cancel != nil {
			d.cancel()
		}
		d.loopWg.Wait()
		d.logger.Info("Block event dequeue synchronizer daemon stopped")
	})
}

func (d *BlockEventsDequeueDaemon) Failed() <-chan error {
	return d.failed
}

func (d *BlockEventsDequeueDaemon) runLoop() {
	defer d.loopWg.Done()
	ticker := time.NewTicker(d.syncSettings.BlockEventsPollInterval)
	defer ticker.Stop()

	for {
		if d.ctx.Err() != nil {
			return
		}

		if err := d.DequeueAndProcess(d.ctx); err != nil {
			if d.ctx.Err() != nil {
				if !errors.Is(err, context.Canceled) {
					d.logger.Error("Block failed during shutdown", zap.Error(err))
				}
				return
			}
			d.failed <- err
			return
		}

		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (d *BlockEventsDequeueDaemon) DequeueAndProcess(ctx context.Context) error {
	blockNumber, item := d.blockEventsQueue.Dequeue(ctx)
	if item == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := d.blockEventsProcessor.ProcessBlock(ctx, blockNumber, item)
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	d.logger.Info("Processed block events successfully",
		zap.Uint64("blockNumber", blockNumber),
	)

	return nil
}
