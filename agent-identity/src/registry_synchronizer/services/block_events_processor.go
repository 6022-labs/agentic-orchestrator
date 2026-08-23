package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/traces/interfaces"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const (
	blockEventsRetryAttempts = 3
	blockEventsRetryDelay    = 100 * time.Millisecond
)

type BlockEventsProcessor struct {
	logger *zap.Logger

	syncMetrics metrics_interfaces.SyncMetricsInterface
	syncTracer  traces_interfaces.SyncTracerInterface

	syncCheckpointMu         sync.Mutex
	syncCheckpointRepository repositories.SyncCheckpointRepositoryInterface

	agentCollectionV1MintedHandler                interfaces.MintedEventHandlerInterface
	agentCollectionV1TransferHandler              interfaces.TransferEventHandlerInterface
	agentCollectionV1AgentImageUpdatedHandler     interfaces.AgentImageUpdatedEventHandlerInterface
	agentCollectionV1AgentImageRemovedHandler     interfaces.AgentImageRemovedEventHandlerInterface
	agentCollectionV1AgentAddressAddedHandler     interfaces.AgentAddressAddedEventHandlerInterface
	agentCollectionV1AgentAddressRemovedHandler   interfaces.AgentAddressRemovedEventHandlerInterface
	agentCollectionV1MintProposalCreatedHandler   interfaces.MintProposalCreatedEventHandlerInterface
	agentCollectionV1MintProposalRemovedHandler   interfaces.MintProposalRemovedEventHandlerInterface
	agentCollectionV1AgentAttributeUpdatedHandler interfaces.AgentAttributeUpdatedEventHandlerInterface

	agentCollectionManagerCollectionCreatedHandler interfaces.CollectionCreatedEventHandlerInterface

	onBlockProcessedHandlers []interfaces.OnBlockProcessedHandlerInterface
}

type NewBlockEventsProcessorParams struct {
	dig.In

	Logger *zap.Logger

	SyncMetrics metrics_interfaces.SyncMetricsInterface
	SyncTracer  traces_interfaces.SyncTracerInterface

	SyncCheckpointRepository repositories.SyncCheckpointRepositoryInterface

	AgentCollectionV1MintedHandler                interfaces.MintedEventHandlerInterface
	AgentCollectionV1TransferHandler              interfaces.TransferEventHandlerInterface
	AgentCollectionV1AgentImageUpdatedHandler     interfaces.AgentImageUpdatedEventHandlerInterface
	AgentCollectionV1AgentImageRemovedHandler     interfaces.AgentImageRemovedEventHandlerInterface
	AgentCollectionV1AgentAddressAddedHandler     interfaces.AgentAddressAddedEventHandlerInterface
	AgentCollectionV1AgentAddressRemovedHandler   interfaces.AgentAddressRemovedEventHandlerInterface
	AgentCollectionV1MintProposalCreatedHandler   interfaces.MintProposalCreatedEventHandlerInterface
	AgentCollectionV1MintProposalRemovedHandler   interfaces.MintProposalRemovedEventHandlerInterface
	AgentCollectionV1AgentAttributeUpdatedHandler interfaces.AgentAttributeUpdatedEventHandlerInterface

	AgentCollectionManagerCollectionCreatedHandler interfaces.CollectionCreatedEventHandlerInterface

	OnBlockProcessedHandlers []interfaces.OnBlockProcessedHandlerInterface `group:"on_block_processed_handlers"`
}

func NewBlockEventsProcessor(
	params NewBlockEventsProcessorParams,
) *BlockEventsProcessor {
	return &BlockEventsProcessor{
		logger: params.Logger,

		syncMetrics: params.SyncMetrics,
		syncTracer:  params.SyncTracer,

		syncCheckpointRepository: params.SyncCheckpointRepository,

		agentCollectionV1MintedHandler:                params.AgentCollectionV1MintedHandler,
		agentCollectionV1TransferHandler:              params.AgentCollectionV1TransferHandler,
		agentCollectionV1AgentImageUpdatedHandler:     params.AgentCollectionV1AgentImageUpdatedHandler,
		agentCollectionV1AgentImageRemovedHandler:     params.AgentCollectionV1AgentImageRemovedHandler,
		agentCollectionV1AgentAddressAddedHandler:     params.AgentCollectionV1AgentAddressAddedHandler,
		agentCollectionV1AgentAddressRemovedHandler:   params.AgentCollectionV1AgentAddressRemovedHandler,
		agentCollectionV1MintProposalCreatedHandler:   params.AgentCollectionV1MintProposalCreatedHandler,
		agentCollectionV1MintProposalRemovedHandler:   params.AgentCollectionV1MintProposalRemovedHandler,
		agentCollectionV1AgentAttributeUpdatedHandler: params.AgentCollectionV1AgentAttributeUpdatedHandler,

		agentCollectionManagerCollectionCreatedHandler: params.AgentCollectionManagerCollectionCreatedHandler,

		onBlockProcessedHandlers: params.OnBlockProcessedHandlers,
	}
}

func (p *BlockEventsProcessor) ProcessBlock(
	ctx context.Context,
	blockNumber uint64,
	item *interfaces.QueueItem,
) (err error) {
	ctx, span := p.syncTracer.StartBlockProcessing(ctx, blockNumber)
	start := time.Now()
	defer func() {
		outcome := metrics_interfaces.BlockOutcomeProcessed
		if err != nil {
			outcome = metrics_interfaces.BlockOutcomeFailed
			span.Fail(err)
		}
		p.syncMetrics.RecordBlockProcessed(ctx, outcome, time.Since(start))
		span.End()
	}()

	p.logger.Info(
		"Processing block events",
		zap.Uint64("blockNumber", blockNumber),
		zap.Int("agentCollectionV1MintedEventsCount", len(item.AgentCollectionV1MintedEvents)),
		zap.Int("agentCollectionV1TransferEventsCount", len(item.AgentCollectionV1TransferEvents)),
		zap.Int("agentCollectionV1AgentAddressAddedEventsCount", len(item.AgentCollectionV1AgentAddressAddedEvents)),
		zap.Int("agentCollectionV1AgentImageUpdatedEventsCount", len(item.AgentCollectionV1AgentImageUpdatedEvents)),
		zap.Int("agentCollectionV1AgentImageRemovedEventsCount", len(item.AgentCollectionV1AgentImageRemovedEvents)),
		zap.Int(
			"agentCollectionV1AgentAddressRemovedEventsCount",
			len(item.AgentCollectionV1AgentAddressRemovedEvents),
		),
		zap.Int(
			"agentCollectionV1MintProposalCreatedEventsCount",
			len(item.AgentCollectionV1MintProposalCreatedEvents),
		),
		zap.Int(
			"agentCollectionV1MintProposalRemovedEventsCount",
			len(item.AgentCollectionV1MintProposalRemovedEvents),
		),
		zap.Int(
			"agentCollectionV1AgentAttributeUpdatedEventsCount",
			len(item.AgentCollectionV1AgentAttributeUpdatedEvents),
		),
		zap.Int(
			"agentCollectionManagerCollectionCreatedEventsCount",
			len(item.AgentCollectionsManagerCollectionCreatedEvents),
		),
	)

	operation := func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := p.syncAgentCollectionsManagerCollectionCreatedEvents(
			ctx,
			item.AgentCollectionsManagerCollectionCreatedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync collection created events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1MintedEvents(ctx, item.AgentCollectionV1MintedEvents); err != nil {
			return fmt.Errorf("failed to sync minted events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1TransferEvents(ctx, item.AgentCollectionV1TransferEvents); err != nil {
			return fmt.Errorf("failed to sync transfer events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentImageUpdatedEvents(
			ctx,
			item.AgentCollectionV1AgentImageUpdatedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync agent image updated events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentImageRemovedEvents(
			ctx,
			item.AgentCollectionV1AgentImageRemovedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync agent image removed events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1MintProposalCreatedEvents(
			ctx,
			item.AgentCollectionV1MintProposalCreatedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync mint proposal created events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1MintProposalRemovedEvents(
			ctx,
			item.AgentCollectionV1MintProposalRemovedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync mint proposal removed events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentAttributeUpdatedEvents(
			ctx,
			item.AgentCollectionV1AgentAttributeUpdatedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync agent attribute updated events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentAddressAddedEvents(
			ctx,
			item.AgentCollectionV1AgentAddressAddedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync agent wallet added events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentAddressRemovedEvents(
			ctx,
			item.AgentCollectionV1AgentAddressRemovedEvents,
		); err != nil {
			return fmt.Errorf("failed to sync agent wallet removed events at block %d: %w", blockNumber, err)
		}

		if err := p.saveOrUpdateSyncCheckpoint(ctx, blockNumber); err != nil {
			return err
		}

		return nil
	}

	if err = retry.Do(
		operation,
		retry.Attempts(blockEventsRetryAttempts),
		retry.Delay(blockEventsRetryDelay),
		retry.DelayType(retry.FixedDelay),
		retry.Context(ctx),
	); err != nil {
		return fmt.Errorf("block %d: failed to process block events: %w", blockNumber, err)
	}

	for _, handler := range p.onBlockProcessedHandlers {
		handler.OnBlockProcessed(ctx, blockNumber)
	}

	p.logger.Info("Processed block events successfully",
		zap.Uint64("blockNumber", blockNumber),
	)

	return nil
}

func (p *BlockEventsProcessor) saveOrUpdateSyncCheckpoint(ctx context.Context, blockNumber uint64) error {
	p.syncCheckpointMu.Lock()
	defer p.syncCheckpointMu.Unlock()

	checkpoint, err := p.syncCheckpointRepository.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sync checkpoint: %w", err)
	}

	if checkpoint == nil {
		newCheckpoint := &models.SyncCheckpoint{ID: uuid.New(), Block: blockNumber}
		if _, err := p.syncCheckpointRepository.Save(ctx, newCheckpoint); err != nil {
			return fmt.Errorf("failed to create sync checkpoint at block %d: %w", blockNumber, err)
		}
		p.syncMetrics.RecordCheckpoint(ctx, blockNumber)
		return nil
	}

	if blockNumber <= checkpoint.Block {
		return nil
	}

	checkpoint.Block = blockNumber
	if _, err := p.syncCheckpointRepository.Save(ctx, checkpoint); err != nil {
		return fmt.Errorf("failed to update sync checkpoint at block %d: %w", blockNumber, err)
	}
	p.syncMetrics.RecordCheckpoint(ctx, blockNumber)
	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1MintedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1Minted,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "Minted")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1MintedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1TransferEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1Transfer,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "Transfer")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1TransferHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentImageUpdatedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1AgentImageUpdated,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "AgentImageUpdated")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1AgentImageUpdatedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentImageRemovedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1AgentImageRemoved,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "AgentImageRemoved")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1AgentImageRemovedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1MintProposalCreatedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1MintProposalCreated,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "MintProposalCreated")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1MintProposalCreatedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1MintProposalRemovedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1MintProposalRemoved,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "MintProposalRemoved")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1MintProposalRemovedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentAttributeUpdatedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1AgentAttributeUpdated,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "AgentAttributeUpdated")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1AgentAttributeUpdatedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentAddressAddedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1AgentAddressAdded,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "AgentAddressAdded")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1AgentAddressAddedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentAddressRemovedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionV1AgentAddressRemoved,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "AgentAddressRemoved")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionV1AgentAddressRemovedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionsManagerCollectionCreatedEvents(
	ctx context.Context,
	events []*abi.AgentCollectionsManagerCollectionCreated,
) error {
	if len(events) == 0 {
		return nil
	}

	ctx, span := p.syncTracer.StartEventHandling(ctx, "CollectionCreated")
	defer span.End()

	for _, event := range events {
		if err := p.agentCollectionManagerCollectionCreatedHandler.Handle(ctx, event); err != nil {
			span.Fail(err)
			return err
		}
	}

	return nil
}
