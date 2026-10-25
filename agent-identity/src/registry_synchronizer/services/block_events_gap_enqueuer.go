package services

import (
	"context"

	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	common_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/metrics/interfaces"
	synchronizer_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type BlockEventsGapEnqueuer struct {
	logger *zap.Logger

	syncMetrics metrics_interfaces.SyncMetricsInterface

	syncSettings *settings.SyncSettings

	blockEventsQueue interfaces.BlockEventsQueueInterface

	blockchainRequester                        interfaces.BlockchainRequesterInterface
	agentCollectionBlockchainRequester         interfaces.AgentCollectionBlockchainRequesterInterface
	agentCollectionsManagerBlockchainRequester interfaces.AgentCollectionsManagerBlockchainRequesterInterface

	agentCollectionRepository common_repositories.AgentCollectionRepositoryInterface
	syncCheckpointRepository  synchronizer_repositories.SyncCheckpointRepositoryInterface
}

func NewBlockEventsGapEnqueuer(
	logger *zap.Logger,

	syncMetrics metrics_interfaces.SyncMetricsInterface,

	syncSettings *settings.SyncSettings,

	blockEventsQueue interfaces.BlockEventsQueueInterface,

	blockchainRequester interfaces.BlockchainRequesterInterface,
	agentCollectionBlockchainRequester interfaces.AgentCollectionBlockchainRequesterInterface,
	agentCollectionsManagerBlockchainRequester interfaces.AgentCollectionsManagerBlockchainRequesterInterface,

	agentCollectionRepository common_repositories.AgentCollectionRepositoryInterface,
	syncCheckpointRepository synchronizer_repositories.SyncCheckpointRepositoryInterface,
) *BlockEventsGapEnqueuer {
	return &BlockEventsGapEnqueuer{
		logger: logger,

		syncMetrics: syncMetrics,

		syncSettings: syncSettings,

		blockEventsQueue: blockEventsQueue,

		blockchainRequester:                        blockchainRequester,
		agentCollectionBlockchainRequester:         agentCollectionBlockchainRequester,
		agentCollectionsManagerBlockchainRequester: agentCollectionsManagerBlockchainRequester,

		syncCheckpointRepository:  syncCheckpointRepository,
		agentCollectionRepository: agentCollectionRepository,
	}
}

func (b *BlockEventsGapEnqueuer) EnqueueGap(ctx context.Context) error {
	toBlock, err := b.blockchainRequester.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("error getting current block number: %w", err)
	}

	syncCheckpoint, err := b.syncCheckpointRepository.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sync checkpoint: %w", err)
	}

	minimumBlockToStartFrom := b.syncSettings.StartBlock
	if syncCheckpoint != nil && syncCheckpoint.Block > minimumBlockToStartFrom {
		minimumBlockToStartFrom = syncCheckpoint.Block + 1
	}

	if *toBlock < minimumBlockToStartFrom {
		b.syncMetrics.RecordBackfilledBlocks(ctx, 0)
		return nil
	}

	b.logger.Info("Enqueuing block events gap",
		zap.Uint64("fromBlock", minimumBlockToStartFrom),
		zap.Uint64("toBlock", *toBlock),
	)

	// Kept across windows: earlier windows' collections aren't persisted yet.
	knownCollections, err := b.loadKnownCollections(ctx)
	if err != nil {
		return err
	}

	windowBlocks := b.syncSettings.GapWindowBlocks
	for windowStart := minimumBlockToStartFrom; ; windowStart += windowBlocks {
		windowEnd := *toBlock
		if *toBlock-windowStart >= windowBlocks {
			windowEnd = windowStart + windowBlocks - 1
		}

		windowEnqueuedBlocks, err := b.enqueueGapWindow(ctx, windowStart, windowEnd, knownCollections)
		if err != nil {
			return fmt.Errorf("failed to enqueue gap window %d-%d: %w", windowStart, windowEnd, err)
		}
		b.syncMetrics.RecordBackfilledBlocks(ctx, windowEnqueuedBlocks)

		if windowEnd == *toBlock {
			break
		}
	}

	return nil
}

func (b *BlockEventsGapEnqueuer) enqueueGapWindow(
	ctx context.Context,
	from, to uint64,
	knownCollections map[common.Address]struct{},
) (int64, error) {
	agentCollectionManagerCollectionCreatedEvents, err := b.getAgentCollectionManagerCollectionCreated(ctx, from, to)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collections manager collection created events: %w", err)
	}

	for _, event := range agentCollectionManagerCollectionCreatedEvents {
		knownCollections[event.CollectionAddress] = struct{}{}
	}
	agentCollectionRelatedAddresses := lo.Keys(knownCollections)

	mintedEvents, err := b.getAgentCollectionV1Minted(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 minted events: %w", err)
	}

	transferEvents, err := b.getAgentCollectionV1Transfer(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 transfer events: %w", err)
	}

	imageUpdatedEvents, err := b.getAgentCollectionV1AgentImageUpdated(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 agent image updated events: %w", err)
	}

	imageRemovedEvents, err := b.getAgentCollectionV1AgentImageRemoved(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 agent image removed events: %w", err)
	}

	attributeUpdatedEvents, err := b.getAgentCollectionV1AgentAttributeUpdated(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 agent attribute updated events: %w", err)
	}

	walletAddedEvents, err := b.getAgentCollectionV1AgentAddressAdded(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 agent wallet added events: %w", err)
	}

	walletRemovedEvents, err := b.getAgentCollectionV1AgentAddressRemoved(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 agent wallet removed events: %w", err)
	}

	mintProposalCreatedEvents, err := b.getAgentCollectionV1MintProposalCreated(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 mint proposal created events: %w", err)
	}

	mintProposalRemovedEvents, err := b.getAgentCollectionV1MintProposalRemoved(ctx,
		from,
		to,
		agentCollectionRelatedAddresses,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent collection v1 mint proposal removed events: %w", err)
	}

	collectionCreatedEventsByBlock := lo.GroupBy(
		agentCollectionManagerCollectionCreatedEvents,
		func(item *abi.AgentCollectionsManagerCollectionCreated) uint64 {
			return item.Raw.BlockNumber
		},
	)

	mintedEventsByBlock := lo.GroupBy(mintedEvents, func(item *abi.AgentCollectionV1Minted) uint64 {
		return item.Raw.BlockNumber
	})

	transferEventsByBlock := lo.GroupBy(transferEvents, func(item *abi.AgentCollectionV1Transfer) uint64 {
		return item.Raw.BlockNumber
	})

	imageUpdatedEventsByBlock := lo.GroupBy(
		imageUpdatedEvents,
		func(item *abi.AgentCollectionV1AgentImageUpdated) uint64 {
			return item.Raw.BlockNumber
		},
	)

	imageRemovedEventsByBlock := lo.GroupBy(
		imageRemovedEvents,
		func(item *abi.AgentCollectionV1AgentImageRemoved) uint64 {
			return item.Raw.BlockNumber
		},
	)

	attributeUpdatedEventsByBlock := lo.GroupBy(
		attributeUpdatedEvents,
		func(item *abi.AgentCollectionV1AgentAttributeUpdated) uint64 {
			return item.Raw.BlockNumber
		},
	)

	walletAddedEventsByBlock := lo.GroupBy(
		walletAddedEvents,
		func(item *abi.AgentCollectionV1AgentAddressAdded) uint64 {
			return item.Raw.BlockNumber
		},
	)

	walletRemovedEventsByBlock := lo.GroupBy(
		walletRemovedEvents,
		func(item *abi.AgentCollectionV1AgentAddressRemoved) uint64 {
			return item.Raw.BlockNumber
		},
	)

	mintProposalCreatedEventsByBlock := lo.GroupBy(
		mintProposalCreatedEvents,
		func(item *abi.AgentCollectionV1MintProposalCreated) uint64 {
			return item.Raw.BlockNumber
		},
	)

	mintProposalRemovedEventsByBlock := lo.GroupBy(
		mintProposalRemovedEvents,
		func(item *abi.AgentCollectionV1MintProposalRemoved) uint64 {
			return item.Raw.BlockNumber
		},
	)

	enqueuedBlocks := int64(0)

	for blockNumber := from; blockNumber <= to; blockNumber++ {
		collectionCreatedEventsForBlock := collectionCreatedEventsByBlock[blockNumber]

		mintedEventsForBlock := mintedEventsByBlock[blockNumber]
		transferEventsForBlock := transferEventsByBlock[blockNumber]
		imageUpdatedEventsForBlock := imageUpdatedEventsByBlock[blockNumber]
		imageRemovedEventsForBlock := imageRemovedEventsByBlock[blockNumber]
		attributeUpdatedEventsForBlock := attributeUpdatedEventsByBlock[blockNumber]
		walletAddedEventsForBlock := walletAddedEventsByBlock[blockNumber]
		walletRemovedEventsForBlock := walletRemovedEventsByBlock[blockNumber]
		mintProposalCreatedEventsForBlock := mintProposalCreatedEventsByBlock[blockNumber]
		mintProposalRemovedEventsForBlock := mintProposalRemovedEventsByBlock[blockNumber]

		if len(collectionCreatedEventsForBlock) == 0 &&
			len(mintedEventsForBlock) == 0 &&
			len(transferEventsForBlock) == 0 &&
			len(imageUpdatedEventsForBlock) == 0 &&
			len(imageRemovedEventsForBlock) == 0 &&
			len(attributeUpdatedEventsForBlock) == 0 &&
			len(walletAddedEventsForBlock) == 0 &&
			len(walletRemovedEventsForBlock) == 0 &&
			len(mintProposalCreatedEventsForBlock) == 0 &&
			len(mintProposalRemovedEventsForBlock) == 0 {
			continue
		}

		enqueuedBlocks++

		if err = b.blockEventsQueue.AddAgentCollectionsManagerCollectionCreatedEvents(
			ctx,
			blockNumber,
			collectionCreatedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collections manager collection created events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1MintedEvents(
			ctx,
			blockNumber,
			mintedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 minted events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1TransferEvents(
			ctx,
			blockNumber,
			transferEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 transfer events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentImageUpdatedEvents(
			ctx,
			blockNumber,
			imageUpdatedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 agent image updated events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentImageRemovedEvents(
			ctx,
			blockNumber,
			imageRemovedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 agent image removed events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentAttributeUpdatedEvents(
			ctx,
			blockNumber,
			attributeUpdatedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 agent attribute updated events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentAddressAddedEvents(
			ctx,
			blockNumber,
			walletAddedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 agent wallet added events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentAddressRemovedEvents(
			ctx,
			blockNumber,
			walletRemovedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 agent wallet removed events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1MintProposalCreatedEvents(
			ctx,
			blockNumber,
			mintProposalCreatedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 mint proposal created events for block %d: %w",
				blockNumber,
				err,
			)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1MintProposalRemovedEvents(
			ctx,
			blockNumber,
			mintProposalRemovedEventsForBlock,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to enqueue agent collection v1 mint proposal removed events for block %d: %w",
				blockNumber,
				err,
			)
		}
	}

	return enqueuedBlocks, nil
}

func (s *BlockEventsGapEnqueuer) loadKnownCollections(ctx context.Context) (map[common.Address]struct{}, error) {
	inDatabaseCollections, err := s.agentCollectionRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collections from database: %w", err)
	}

	knownCollections := make(map[common.Address]struct{}, len(inDatabaseCollections))
	for _, collection := range inDatabaseCollections {
		knownCollections[common.HexToAddress(collection.Address)] = struct{}{}
	}

	return knownCollections, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1Minted(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1Minted, error) {
	s.logger.Debug("Starting to get agent collection v1 minted events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1MintedEvent(
		ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 minted events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 minted events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1Transfer(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1Transfer, error) {
	s.logger.Debug("Starting to get agent collection v1 transfer events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1TransferEvent(
		ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 transfer events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 transfer events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentImageUpdated(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1AgentImageUpdated, error) {
	s.logger.Debug("Starting to get agent collection v1 agent image updated events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentImageUpdatedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent image updated events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent image updated events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentImageRemoved(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1AgentImageRemoved, error) {
	s.logger.Debug("Starting to get agent collection v1 agent image removed events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentImageRemovedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent image removed events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent image removed events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentAttributeUpdated(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1AgentAttributeUpdated, error) {
	s.logger.Debug("Starting to get agent collection v1 agent attribute updated events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentAttributeUpdatedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent attribute updated events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent attribute updated events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentAddressAdded(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1AgentAddressAdded, error) {
	s.logger.Debug("Starting to get agent collection v1 agent wallet added events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentAddressAddedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent wallet added events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent wallet added events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentAddressRemoved(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1AgentAddressRemoved, error) {
	s.logger.Debug("Starting to get agent collection v1 agent wallet removed events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentAddressRemovedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent wallet removed events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent wallet removed events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1MintProposalCreated(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1MintProposalCreated, error) {
	s.logger.Debug("Starting to get agent collection v1 mint proposal created events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1MintProposalCreatedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 mint proposal created events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 mint proposal created events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1MintProposalRemoved(
	ctx context.Context,
	from, to uint64,
	collectionAddresses []common.Address,
) ([]*abi.AgentCollectionV1MintProposalRemoved, error) {
	s.logger.Debug("Starting to get agent collection v1 mint proposal removed events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1MintProposalRemovedEvent(ctx,
		from,
		to,
		collectionAddresses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 mint proposal removed events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 mint proposal removed events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionManagerCollectionCreated(
	ctx context.Context,
	from, to uint64,
) ([]*abi.AgentCollectionsManagerCollectionCreated, error) {
	s.logger.Debug("Starting to get agent collections manager collection created events")

	events, err := s.agentCollectionsManagerBlockchainRequester.GetAgentCollectionsManagerCollectionCreatedEvent(ctx,
		from,
		to,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collections manager collection created events: %w", err)
	}

	s.logger.Debug("Fetched agent collections manager collection created events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}
