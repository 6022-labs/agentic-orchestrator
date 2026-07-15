package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	blockEventsRetryAttempts = 3
	blockEventsRetryDelay    = 100 * time.Millisecond
)

type BlockEventsProcessor struct {
	logger *zap.Logger

	syncCheckpointMu         sync.Mutex
	syncCheckpointRepository repositories.SyncCheckpointRepositoryInterface

	agentCollectionV1MintedHandler                interfaces.MintedEventHandlerInterface
	agentCollectionV1TransferHandler              interfaces.TransferEventHandlerInterface
	agentCollectionV1AgentImageUpdatedHandler     interfaces.AgentImageUpdatedEventHandlerInterface
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

	SyncCheckpointRepository repositories.SyncCheckpointRepositoryInterface

	AgentCollectionV1MintedHandler                interfaces.MintedEventHandlerInterface
	AgentCollectionV1TransferHandler              interfaces.TransferEventHandlerInterface
	AgentCollectionV1AgentImageUpdatedHandler     interfaces.AgentImageUpdatedEventHandlerInterface
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

		syncCheckpointRepository: params.SyncCheckpointRepository,

		agentCollectionV1MintedHandler:                params.AgentCollectionV1MintedHandler,
		agentCollectionV1TransferHandler:              params.AgentCollectionV1TransferHandler,
		agentCollectionV1AgentImageUpdatedHandler:     params.AgentCollectionV1AgentImageUpdatedHandler,
		agentCollectionV1AgentAddressAddedHandler:     params.AgentCollectionV1AgentAddressAddedHandler,
		agentCollectionV1AgentAddressRemovedHandler:   params.AgentCollectionV1AgentAddressRemovedHandler,
		agentCollectionV1MintProposalCreatedHandler:   params.AgentCollectionV1MintProposalCreatedHandler,
		agentCollectionV1MintProposalRemovedHandler:   params.AgentCollectionV1MintProposalRemovedHandler,
		agentCollectionV1AgentAttributeUpdatedHandler: params.AgentCollectionV1AgentAttributeUpdatedHandler,

		agentCollectionManagerCollectionCreatedHandler: params.AgentCollectionManagerCollectionCreatedHandler,

		onBlockProcessedHandlers: params.OnBlockProcessedHandlers,
	}
}

func (p *BlockEventsProcessor) ProcessBlock(ctx context.Context, blockNumber uint64, item *interfaces.QueueItem) error {
	p.logger.Info("Processing block events",
		zap.Uint64("blockNumber", blockNumber),
		zap.Int("agentCollectionV1MintedEventsCount", len(item.AgentCollectionV1MintedEvents)),
		zap.Int("agentCollectionV1TransferEventsCount", len(item.AgentCollectionV1TransferEvents)),
		zap.Int("agentCollectionV1AgentAddressAddedEventsCount", len(item.AgentCollectionV1AgentAddressAddedEvents)),
		zap.Int("agentCollectionV1AgentImageUpdatedEventsCount", len(item.AgentCollectionV1AgentImageUpdatedEvents)),
		zap.Int("agentCollectionV1AgentAddressRemovedEventsCount", len(item.AgentCollectionV1AgentAddressRemovedEvents)),
		zap.Int("agentCollectionV1MintProposalCreatedEventsCount", len(item.AgentCollectionV1MintProposalCreatedEvents)),
		zap.Int("agentCollectionV1MintProposalRemovedEventsCount", len(item.AgentCollectionV1MintProposalRemovedEvents)),
		zap.Int("agentCollectionV1AgentAttributeUpdatedEventsCount", len(item.AgentCollectionV1AgentAttributeUpdatedEvents)),
		zap.Int("agentCollectionManagerCollectionCreatedEventsCount", len(item.AgentCollectionsManagerCollectionCreatedEvents)),
	)

	operation := func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := p.syncAgentCollectionsManagerCollectionCreatedEvents(item.AgentCollectionsManagerCollectionCreatedEvents); err != nil {
			return fmt.Errorf("failed to sync collection created events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1MintedEvents(item.AgentCollectionV1MintedEvents); err != nil {
			return fmt.Errorf("failed to sync minted events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1TransferEvents(item.AgentCollectionV1TransferEvents); err != nil {
			return fmt.Errorf("failed to sync transfer events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentImageUpdatedEvents(item.AgentCollectionV1AgentImageUpdatedEvents); err != nil {
			return fmt.Errorf("failed to sync agent image updated events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1MintProposalCreatedEvents(item.AgentCollectionV1MintProposalCreatedEvents); err != nil {
			return fmt.Errorf("failed to sync mint proposal created events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1MintProposalRemovedEvents(item.AgentCollectionV1MintProposalRemovedEvents); err != nil {
			return fmt.Errorf("failed to sync mint proposal removed events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentAttributeUpdatedEvents(item.AgentCollectionV1AgentAttributeUpdatedEvents); err != nil {
			return fmt.Errorf("failed to sync agent attribute updated events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentAddressAddedEvents(item.AgentCollectionV1AgentAddressAddedEvents); err != nil {
			return fmt.Errorf("failed to sync agent wallet added events at block %d: %w", blockNumber, err)
		}

		if err := p.syncAgentCollectionV1AgentAddressRemovedEvents(item.AgentCollectionV1AgentAddressRemovedEvents); err != nil {
			return fmt.Errorf("failed to sync agent wallet removed events at block %d: %w", blockNumber, err)
		}

		if err := p.saveOrUpdateSyncCheckpoint(blockNumber); err != nil {
			return err
		}

		return nil
	}

	if err := retry.Do(
		operation,
		retry.Attempts(blockEventsRetryAttempts),
		retry.Delay(blockEventsRetryDelay),
		retry.DelayType(retry.FixedDelay),
		retry.Context(ctx),
	); err != nil {
		return fmt.Errorf("block %d: failed to process block events: %w", blockNumber, err)
	}

	for _, handler := range p.onBlockProcessedHandlers {
		handler.OnBlockProcessed(blockNumber)
	}

	p.logger.Info("Processed block events successfully",
		zap.Uint64("blockNumber", blockNumber),
	)

	return nil
}

func (p *BlockEventsProcessor) saveOrUpdateSyncCheckpoint(blockNumber uint64) error {
	p.syncCheckpointMu.Lock()
	defer p.syncCheckpointMu.Unlock()

	checkpoint, err := p.syncCheckpointRepository.Get()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get sync checkpoint: %w", err)
	}

	if checkpoint == nil {
		newCheckpoint := &models.SyncCheckpoint{ID: uuid.New(), Block: blockNumber}
		if _, err := p.syncCheckpointRepository.Save(newCheckpoint); err != nil {
			return fmt.Errorf("failed to create sync checkpoint at block %d: %w", blockNumber, err)
		}
		return nil
	}

	if blockNumber <= checkpoint.Block {
		return nil
	}

	checkpoint.Block = blockNumber
	if _, err := p.syncCheckpointRepository.Save(checkpoint); err != nil {
		return fmt.Errorf("failed to update sync checkpoint at block %d: %w", blockNumber, err)
	}
	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1MintedEvents(events []*abi.AgentCollectionV1Minted) error {
	for _, event := range events {
		if err := p.agentCollectionV1MintedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1TransferEvents(events []*abi.AgentCollectionV1Transfer) error {
	for _, event := range events {
		if err := p.agentCollectionV1TransferHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentImageUpdatedEvents(events []*abi.AgentCollectionV1AgentImageUpdated) error {
	for _, event := range events {
		if err := p.agentCollectionV1AgentImageUpdatedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1MintProposalCreatedEvents(events []*abi.AgentCollectionV1MintProposalCreated) error {
	for _, event := range events {
		if err := p.agentCollectionV1MintProposalCreatedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1MintProposalRemovedEvents(events []*abi.AgentCollectionV1MintProposalRemoved) error {
	for _, event := range events {
		if err := p.agentCollectionV1MintProposalRemovedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentAttributeUpdatedEvents(events []*abi.AgentCollectionV1AgentAttributeUpdated) error {
	for _, event := range events {
		if err := p.agentCollectionV1AgentAttributeUpdatedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentAddressAddedEvents(events []*abi.AgentCollectionV1AgentAddressAdded) error {
	for _, event := range events {
		if err := p.agentCollectionV1AgentAddressAddedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentAddressRemovedEvents(events []*abi.AgentCollectionV1AgentAddressRemoved) error {
	for _, event := range events {
		if err := p.agentCollectionV1AgentAddressRemovedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *BlockEventsProcessor) syncAgentCollectionsManagerCollectionCreatedEvents(events []*abi.AgentCollectionsManagerCollectionCreated) error {
	for _, event := range events {
		if err := p.agentCollectionManagerCollectionCreatedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}
