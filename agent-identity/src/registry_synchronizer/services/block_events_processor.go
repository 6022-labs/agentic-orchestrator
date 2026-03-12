package services

import (
	"context"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/avast/retry-go/v4"
	"go.uber.org/zap"
)

const (
	blockEventsRetryAttempts = 3
	blockEventsRetryDelay    = 100 * time.Millisecond
)

type BlockEventsProcessor struct {
	logger *zap.Logger

	syncCheckpointRepository repositories.SyncCheckpointRepositoryInterface

	agentCollectionV1MintedHandler                interfaces.MintedEventHandlerInterface
	agentCollectionV1TransferHandler              interfaces.TransferEventHandlerInterface
	agentCollectionV1AgentImageUpdatedHandler     interfaces.AgentImageUpdatedEventHandlerInterface
	agentCollectionV1AgentSignerUpdatedHandler    interfaces.AgentSignerUpdatedEventHandlerInterface
	agentCollectionV1MintProposalCreatedHandler   interfaces.MintProposalCreatedEventHandlerInterface
	agentCollectionV1MintProposalRemovedHandler   interfaces.MintProposalRemovedEventHandlerInterface
	agentCollectionV1AgentAttributeUpdatedHandler interfaces.AgentAttributeUpdatedEventHandlerInterface

	agentCollectionManagerCollectionCreatedHandler interfaces.CollectionCreatedEventHandlerInterface
}

func NewBlockEventsProcessor(
	logger *zap.Logger,

	syncCheckpointRepository repositories.SyncCheckpointRepositoryInterface,

	agentCollectionV1MintedHandler interfaces.MintedEventHandlerInterface,
	agentCollectionV1TransferHandler interfaces.TransferEventHandlerInterface,
	agentCollectionV1AgentImageUpdatedHandler interfaces.AgentImageUpdatedEventHandlerInterface,
	agentCollectionV1AgentSignerUpdatedHandler interfaces.AgentSignerUpdatedEventHandlerInterface,
	agentCollectionV1MintProposalCreatedHandler interfaces.MintProposalCreatedEventHandlerInterface,
	agentCollectionV1MintProposalRemovedHandler interfaces.MintProposalRemovedEventHandlerInterface,
	agentCollectionV1AgentAttributeUpdatedHandler interfaces.AgentAttributeUpdatedEventHandlerInterface,

	agentCollectionManagerCollectionCreatedHandler interfaces.CollectionCreatedEventHandlerInterface,
) *BlockEventsProcessor {
	return &BlockEventsProcessor{
		logger: logger,

		syncCheckpointRepository: syncCheckpointRepository,

		agentCollectionV1MintedHandler:                agentCollectionV1MintedHandler,
		agentCollectionV1TransferHandler:              agentCollectionV1TransferHandler,
		agentCollectionV1AgentImageUpdatedHandler:     agentCollectionV1AgentImageUpdatedHandler,
		agentCollectionV1AgentSignerUpdatedHandler:    agentCollectionV1AgentSignerUpdatedHandler,
		agentCollectionV1MintProposalCreatedHandler:   agentCollectionV1MintProposalCreatedHandler,
		agentCollectionV1MintProposalRemovedHandler:   agentCollectionV1MintProposalRemovedHandler,
		agentCollectionV1AgentAttributeUpdatedHandler: agentCollectionV1AgentAttributeUpdatedHandler,

		agentCollectionManagerCollectionCreatedHandler: agentCollectionManagerCollectionCreatedHandler,
	}
}

func (p *BlockEventsProcessor) ProcessBlock(ctx context.Context, blockNumber uint64, item *interfaces.QueueItem) error {
	p.logger.Info("Processing block events",
		zap.Uint64("blockNumber", blockNumber),
		zap.Int("agentCollectionV1MintedEventsCount", len(item.AgentCollectionV1MintedEvents)),
		zap.Int("agentCollectionV1TransferEventsCount", len(item.AgentCollectionV1TransferEvents)),
		zap.Int("agentCollectionV1AgentImageUpdatedEventsCount", len(item.AgentCollectionV1AgentImageUpdatedEvents)),
		zap.Int("agentCollectionV1AgentSignerUpdatedEventsCount", len(item.AgentCollectionV1AgentSignerUpdatedEvents)),
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
			return fmt.Errorf("failed to sync collection created events: %w", err)
		}

		if err := p.syncAgentCollectionV1MintedEvents(item.AgentCollectionV1MintedEvents); err != nil {
			return fmt.Errorf("failed to sync minted events: %w", err)
		}

		if err := p.syncAgentCollectionV1TransferEvents(item.AgentCollectionV1TransferEvents); err != nil {
			return fmt.Errorf("failed to sync transfer events: %w", err)
		}

		if err := p.syncAgentCollectionV1AgentImageUpdatedEvents(item.AgentCollectionV1AgentImageUpdatedEvents); err != nil {
			return fmt.Errorf("failed to sync agent image updated events: %w", err)
		}

		if err := p.syncAgentCollectionV1AgentSignerUpdatedEvents(item.AgentCollectionV1AgentSignerUpdatedEvents); err != nil {
			return fmt.Errorf("failed to sync agent signer updated events: %w", err)
		}

		if err := p.syncAgentCollectionV1MintProposalCreatedEvents(item.AgentCollectionV1MintProposalCreatedEvents); err != nil {
			return fmt.Errorf("failed to sync mint proposal created events: %w", err)
		}

		if err := p.syncAgentCollectionV1MintProposalRemovedEvents(item.AgentCollectionV1MintProposalRemovedEvents); err != nil {
			return fmt.Errorf("failed to sync mint proposal removed events: %w", err)
		}

		if err := p.syncAgentCollectionV1AgentAttributeUpdatedEvents(item.AgentCollectionV1AgentAttributeUpdatedEvents); err != nil {
			return fmt.Errorf("failed to sync agent attribute updated events: %w", err)
		}

		syncCheckpoint := &models.SyncCheckpoint{
			Block: blockNumber,
		}

		if _, err := p.syncCheckpointRepository.Save(syncCheckpoint); err != nil {
			return fmt.Errorf("failed to save sync checkpoint at block %d: %w", blockNumber, err)
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

	p.logger.Info("Processed block events successfully",
		zap.Uint64("blockNumber", blockNumber),
	)

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

func (p *BlockEventsProcessor) syncAgentCollectionV1AgentSignerUpdatedEvents(events []*abi.AgentCollectionV1AgentSignerUpdated) error {
	for _, event := range events {
		if err := p.agentCollectionV1AgentSignerUpdatedHandler.Handle(event); err != nil {
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

func (p *BlockEventsProcessor) syncAgentCollectionsManagerCollectionCreatedEvents(events []*abi.AgentCollectionsManagerCollectionCreated) error {
	for _, event := range events {
		if err := p.agentCollectionManagerCollectionCreatedHandler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}
