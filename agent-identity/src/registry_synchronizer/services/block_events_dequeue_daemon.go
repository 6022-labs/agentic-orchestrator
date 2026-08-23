package services

import (
	"context"
	"errors"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/zap"
)

const (
	defaultBlockEventsPollInterval = 200 * time.Millisecond
)

type BlockEventsDequeueDaemon struct {
	logger *zap.Logger

	blockEventsQueue interfaces.BlockEventsQueueInterface

	blockEventsProcessor interfaces.BlockEventsProcessorInterface

	cancel  context.CancelFunc
	running bool
}

func NewBlockEventsDequeueDaemon(
	logger *zap.Logger,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	blockEventsProcessor interfaces.BlockEventsProcessorInterface,
) *BlockEventsDequeueDaemon {
	return &BlockEventsDequeueDaemon{
		logger:               logger,
		blockEventsQueue:     blockEventsQueue,
		blockEventsProcessor: blockEventsProcessor,
	}
}

func (d *BlockEventsDequeueDaemon) Start(ctx context.Context) error {
	if d.running {
		return errors.New("block events dequeue synchronizer daemon already running")
	}
	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel
	d.running = true

	d.logger.Info("Starting block event dequeue synchronizer daemon...")

	defer func() {
		d.cancel = nil
		d.running = false
	}()

	err := d.run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}

func (d *BlockEventsDequeueDaemon) Stop() error {
	if !d.running || d.cancel == nil {
		return nil
	}
	cancel := d.cancel

	d.logger.Info("Stopping block event dequeue synchronizer daemon...")
	cancel()

	return nil
}

func (d *BlockEventsDequeueDaemon) run(ctx context.Context) error {
	ticker := time.NewTicker(defaultBlockEventsPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := d.DequeueAndProcess(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
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
