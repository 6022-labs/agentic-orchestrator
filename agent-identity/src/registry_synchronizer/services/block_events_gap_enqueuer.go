package services

import (
	"errors"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	common_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	synchronizer_repositories "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BlockEventsGapEnqueuer struct {
	logger *zap.Logger

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

		syncSettings: syncSettings,

		blockEventsQueue: blockEventsQueue,

		blockchainRequester:                        blockchainRequester,
		agentCollectionBlockchainRequester:         agentCollectionBlockchainRequester,
		agentCollectionsManagerBlockchainRequester: agentCollectionsManagerBlockchainRequester,

		syncCheckpointRepository:  syncCheckpointRepository,
		agentCollectionRepository: agentCollectionRepository,
	}
}

func (b *BlockEventsGapEnqueuer) EnqueueGap() error {
	toBlock, err := b.blockchainRequester.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("error getting current block number: %w", err)
	}

	syncCheckpoint, err := b.syncCheckpointRepository.Get()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get sync checkpoint: %w", err)
	}

	minimumBlockToStartFrom := b.syncSettings.StartBlock
	if syncCheckpoint != nil && syncCheckpoint.Block > minimumBlockToStartFrom {
		minimumBlockToStartFrom = syncCheckpoint.Block + 1
	}

	b.logger.Info("Enqueuing block events gap",
		zap.Uint64("fromBlock", minimumBlockToStartFrom),
		zap.Uint64("toBlock", *toBlock),
	)

	agentCollectionManagerCollectionCreatedEvents, err := b.getAgentCollectionManagerCollectionCreated(minimumBlockToStartFrom, *toBlock)
	if err != nil {
		return fmt.Errorf("failed to get agent collections manager collection created events: %w", err)
	}

	agentCollectionRelatedAddresses, err := b.buildAgentCollectionRelatedAddresses(agentCollectionManagerCollectionCreatedEvents)
	if err != nil {
		return fmt.Errorf("failed to build agent collection related addresses: %w", err)
	}

	mintedEvents, err := b.getAgentCollectionV1Minted(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 minted events: %w", err)
	}

	transferEvents, err := b.getAgentCollectionV1Transfer(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 transfer events: %w", err)
	}

	imageUpdatedEvents, err := b.getAgentCollectionV1AgentImageUpdated(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 agent image updated events: %w", err)
	}

	signerUpdatedEvents, err := b.getAgentCollectionV1AgentSignerUpdated(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 agent signer updated events: %w", err)
	}

	attributeUpdatedEvents, err := b.getAgentCollectionV1AgentAttributeUpdated(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 agent attribute updated events: %w", err)
	}

	mintProposalCreatedEvents, err := b.getAgentCollectionV1MintProposalCreated(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 mint proposal created events: %w", err)
	}

	mintProposalRemovedEvents, err := b.getAgentCollectionV1MintProposalRemoved(minimumBlockToStartFrom, *toBlock, agentCollectionRelatedAddresses)
	if err != nil {
		return fmt.Errorf("failed to get agent collection v1 mint proposal removed events: %w", err)
	}

	collectionCreatedEventsByBlock := lo.GroupBy(agentCollectionManagerCollectionCreatedEvents, func(item *abi.AgentCollectionsManagerCollectionCreated) uint64 {
		return item.Raw.BlockNumber
	})

	mintedEventsByBlock := lo.GroupBy(mintedEvents, func(item *abi.AgentCollectionV1Minted) uint64 {
		return item.Raw.BlockNumber
	})

	transferEventsByBlock := lo.GroupBy(transferEvents, func(item *abi.AgentCollectionV1Transfer) uint64 {
		return item.Raw.BlockNumber
	})

	imageUpdatedEventsByBlock := lo.GroupBy(imageUpdatedEvents, func(item *abi.AgentCollectionV1AgentImageUpdated) uint64 {
		return item.Raw.BlockNumber
	})

	signerUpdatedEventsByBlock := lo.GroupBy(signerUpdatedEvents, func(item *abi.AgentCollectionV1AgentSignerUpdated) uint64 {
		return item.Raw.BlockNumber
	})

	attributeUpdatedEventsByBlock := lo.GroupBy(attributeUpdatedEvents, func(item *abi.AgentCollectionV1AgentAttributeUpdated) uint64 {
		return item.Raw.BlockNumber
	})

	mintProposalCreatedEventsByBlock := lo.GroupBy(mintProposalCreatedEvents, func(item *abi.AgentCollectionV1MintProposalCreated) uint64 {
		return item.Raw.BlockNumber
	})

	mintProposalRemovedEventsByBlock := lo.GroupBy(mintProposalRemovedEvents, func(item *abi.AgentCollectionV1MintProposalRemoved) uint64 {
		return item.Raw.BlockNumber
	})

	for blockNumber := minimumBlockToStartFrom; blockNumber <= *toBlock; blockNumber++ {
		collectionCreatedEventsForBlock := collectionCreatedEventsByBlock[blockNumber]

		mintedEventsForBlock := mintedEventsByBlock[blockNumber]
		transferEventsForBlock := transferEventsByBlock[blockNumber]
		imageUpdatedEventsForBlock := imageUpdatedEventsByBlock[blockNumber]
		signerUpdatedEventsForBlock := signerUpdatedEventsByBlock[blockNumber]
		attributeUpdatedEventsForBlock := attributeUpdatedEventsByBlock[blockNumber]
		mintProposalCreatedEventsForBlock := mintProposalCreatedEventsByBlock[blockNumber]
		mintProposalRemovedEventsForBlock := mintProposalRemovedEventsByBlock[blockNumber]

		if len(collectionCreatedEventsForBlock) == 0 &&
			len(mintedEventsForBlock) == 0 &&
			len(transferEventsForBlock) == 0 &&
			len(signerUpdatedEventsForBlock) == 0 &&
			len(imageUpdatedEventsForBlock) == 0 &&
			len(attributeUpdatedEventsForBlock) == 0 &&
			len(mintProposalCreatedEventsForBlock) == 0 &&
			len(mintProposalRemovedEventsForBlock) == 0 {
			continue
		}

		if err = b.blockEventsQueue.AddAgentCollectionsManagerCollectionCreatedEvents(blockNumber, collectionCreatedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collections manager collection created events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1MintedEvents(blockNumber, mintedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 minted events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1TransferEvents(blockNumber, transferEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 transfer events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentImageUpdatedEvents(blockNumber, imageUpdatedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 agent image updated events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentSignerUpdatedEvents(blockNumber, signerUpdatedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 agent signer updated events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1AgentAttributeUpdatedEvents(blockNumber, attributeUpdatedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 agent attribute updated events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1MintProposalCreatedEvents(blockNumber, mintProposalCreatedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 mint proposal created events for block %d: %w", blockNumber, err)
		}

		if err = b.blockEventsQueue.AddAgentCollectionV1MintProposalRemovedEvents(blockNumber, mintProposalRemovedEventsForBlock); err != nil {
			return fmt.Errorf("failed to enqueue agent collection v1 mint proposal removed events for block %d: %w", blockNumber, err)
		}

		b.blockEventsQueue.HandleBlockProcessed(blockNumber)
	}

	return nil
}

func (s *BlockEventsGapEnqueuer) buildAgentCollectionRelatedAddresses(agentCollectionManagerCollectionCreatedEvents []*abi.AgentCollectionsManagerCollectionCreated) ([]common.Address, error) {
	agentCollectionAddresses := lo.Map(agentCollectionManagerCollectionCreatedEvents, func(item *abi.AgentCollectionsManagerCollectionCreated, _ int) common.Address {
		return item.CollectionAddress
	})

	inDatabaseCollections, err := s.agentCollectionRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collections from database: %w", err)
	}

	agentCollectionAddressesFromDatabase := lo.Map(inDatabaseCollections, func(item *models.AgentCollection, _ int) common.Address {
		return common.HexToAddress(item.Address)
	})

	agentCollectionAddresses = lo.Uniq(append(agentCollectionAddresses, agentCollectionAddressesFromDatabase...))

	return agentCollectionAddresses, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1Minted(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1Minted, error) {
	s.logger.Debug("Starting to get agent collection v1 minted events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1MintedEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 minted events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 minted events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1Transfer(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1Transfer, error) {
	s.logger.Debug("Starting to get agent collection v1 transfer events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1TransferEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 transfer events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 transfer events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentSignerUpdated(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1AgentSignerUpdated, error) {
	s.logger.Debug("Starting to get agent collection v1 agent signer updated events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentSignerUpdatedEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent signer updated events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent signer updated events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentImageUpdated(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1AgentImageUpdated, error) {
	s.logger.Debug("Starting to get agent collection v1 agent image updated events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentImageUpdatedEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent image updated events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent image updated events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1AgentAttributeUpdated(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1AgentAttributeUpdated, error) {
	s.logger.Debug("Starting to get agent collection v1 agent attribute updated events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1AgentAttributeUpdatedEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 agent attribute updated events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 agent attribute updated events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1MintProposalCreated(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1MintProposalCreated, error) {
	s.logger.Debug("Starting to get agent collection v1 mint proposal created events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1MintProposalCreatedEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 mint proposal created events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 mint proposal created events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionV1MintProposalRemoved(from, to uint64, collectionAddresses []common.Address) ([]*abi.AgentCollectionV1MintProposalRemoved, error) {
	s.logger.Debug("Starting to get agent collection v1 mint proposal removed events")

	events, err := s.agentCollectionBlockchainRequester.GetAgentCollectionV1MintProposalRemovedEvent(from, to, collectionAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collection v1 mint proposal removed events: %w", err)
	}

	s.logger.Debug("Fetched agent collection v1 mint proposal removed events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}

func (s *BlockEventsGapEnqueuer) getAgentCollectionManagerCollectionCreated(from, to uint64) ([]*abi.AgentCollectionsManagerCollectionCreated, error) {
	s.logger.Debug("Starting to get agent collections manager collection created events")

	events, err := s.agentCollectionsManagerBlockchainRequester.GetAgentCollectionsManagerCollectionCreatedEvent(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent collections manager collection created events: %w", err)
	}

	s.logger.Debug("Fetched agent collections manager collection created events",
		zap.Int("eventsCount", len(events)),
	)

	return events, nil
}
