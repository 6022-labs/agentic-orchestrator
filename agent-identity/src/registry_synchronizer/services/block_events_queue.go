package services

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"go.uber.org/zap"
)

const defaultBlockFinalizationDelay = 2 * time.Second

type BlockEventsQueue struct {
	logger       *zap.Logger
	syncMetrics  metrics_interfaces.SyncMetricsInterface
	syncSettings *settings.SyncSettings

	mu                   sync.Mutex
	entries              map[uint64]*queueEntry
	readyList            []uint64
	pendingBlocks        map[uint64]*queueEntry
	latestCompletedBlock uint64
	latestBlockSeen      uint64
	timer                *time.Timer
}

type queueEntry struct {
	data      *interfaces.QueueItem
	hasEvents bool
	complete  bool
	enqueued  bool
}

func NewBlockEventsQueue(
	logger *zap.Logger,
	syncMetrics metrics_interfaces.SyncMetricsInterface,
	syncSettings *settings.SyncSettings,
) *BlockEventsQueue {
	return &BlockEventsQueue{
		logger:       logger,
		syncMetrics:  syncMetrics,
		syncSettings: syncSettings,

		entries:       make(map[uint64]*queueEntry),
		pendingBlocks: make(map[uint64]*queueEntry),
	}
}

func (s *BlockEventsQueue) AddAgentCollectionV1MintedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1Minted,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintedEvents = append(entry.data.AgentCollectionV1MintedEvents, events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1TransferEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1Transfer,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1TransferEvents = append(entry.data.AgentCollectionV1TransferEvents, events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentImageUpdatedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1AgentImageUpdated,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentImageUpdatedEvents = append(
		entry.data.AgentCollectionV1AgentImageUpdatedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentImageRemovedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1AgentImageRemoved,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentImageRemovedEvents = append(
		entry.data.AgentCollectionV1AgentImageRemovedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1MintProposalCreatedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1MintProposalCreated,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintProposalCreatedEvents = append(
		entry.data.AgentCollectionV1MintProposalCreatedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1MintProposalRemovedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1MintProposalRemoved,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintProposalRemovedEvents = append(
		entry.data.AgentCollectionV1MintProposalRemovedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentAttributeUpdatedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1AgentAttributeUpdated,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentAttributeUpdatedEvents = append(
		entry.data.AgentCollectionV1AgentAttributeUpdatedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentAddressAddedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1AgentAddressAdded,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentAddressAddedEvents = append(
		entry.data.AgentCollectionV1AgentAddressAddedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentAddressRemovedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionV1AgentAddressRemoved,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentAddressRemovedEvents = append(
		entry.data.AgentCollectionV1AgentAddressRemovedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionsManagerCollectionCreatedEvents(
	ctx context.Context,
	blockNumber uint64,
	events []*abi.AgentCollectionsManagerCollectionCreated,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionsManagerCollectionCreatedEvents = append(
		entry.data.AgentCollectionsManagerCollectionCreatedEvents,
		events...)
	s.trackPendingBlock(ctx, blockNumber, entry)
	s.recordQueueState(ctx)

	return nil
}

func (s *BlockEventsQueue) Dequeue(ctx context.Context) (uint64, *interfaces.QueueItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.readyList) == 0 {
		return 0, nil
	}

	blockNumber := s.readyList[0]
	s.readyList = s.readyList[1:]

	entry := s.entries[blockNumber]
	delete(s.entries, blockNumber)
	delete(s.pendingBlocks, blockNumber)

	s.recordQueueState(ctx)

	return blockNumber, entry.data
}

func (s *BlockEventsQueue) ensureEntry(blockNumber uint64) *queueEntry {
	entry, exists := s.entries[blockNumber]
	if !exists {
		entry = &queueEntry{
			data: &interfaces.QueueItem{},
		}
		s.entries[blockNumber] = entry
	}

	return entry
}

func (s *BlockEventsQueue) enqueueIfReady(blockNumber uint64, entry *queueEntry) {
	if entry.enqueued {
		return
	}

	if !entry.complete || !entry.hasEvents {
		return
	}

	// Gap blocks can become ready after higher live blocks.
	index := sort.Search(len(s.readyList), func(i int) bool { return s.readyList[i] > blockNumber })
	s.readyList = slices.Insert(s.readyList, index, blockNumber)
	entry.enqueued = true
	delete(s.pendingBlocks, blockNumber)
}

func (s *BlockEventsQueue) trackPendingBlock(ctx context.Context, blockNumber uint64, entry *queueEntry) {
	entry.hasEvents = true
	if entry.enqueued {
		return
	}
	if blockNumber <= s.latestCompletedBlock {
		entry.complete = true
		s.enqueueIfReady(blockNumber, entry)
		return
	}

	s.pendingBlocks[blockNumber] = entry
	s.scheduleBlockFinalization(ctx, blockNumber)
}

func (s *BlockEventsQueue) scheduleBlockFinalization(ctx context.Context, blockNumber uint64) {
	if blockNumber > s.latestBlockSeen {
		s.latestBlockSeen = blockNumber
		s.syncMetrics.RecordHighestBlockWithEvents(ctx, blockNumber)
	}

	if s.timer != nil {
		return
	}

	// The timer outlives the enqueue that scheduled it, so it carries no caller ctx of its own.
	s.timer = time.AfterFunc(s.getDeliveryBuffer(), func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.timer = nil

		completedBlock := s.latestBlockSeen
		if completedBlock > s.latestCompletedBlock {
			s.latestCompletedBlock = completedBlock
		}

		s.completePendingBlocks()
		s.recordQueueState(context.Background())
	})
}

func (s *BlockEventsQueue) getDeliveryBuffer() time.Duration {
	if s.syncSettings.BlockFinalizationDelay > 0 {
		return s.syncSettings.BlockFinalizationDelay
	}

	return defaultBlockFinalizationDelay
}

func (s *BlockEventsQueue) completePendingBlocks() {
	if len(s.pendingBlocks) == 0 {
		return
	}

	var eligible []uint64
	for blockNumber := range s.pendingBlocks {
		if blockNumber <= s.latestCompletedBlock {
			eligible = append(eligible, blockNumber)
		}
	}

	sort.Slice(eligible, func(i, j int) bool {
		return eligible[i] < eligible[j]
	})

	for _, blockNumber := range eligible {
		entry := s.pendingBlocks[blockNumber]
		if entry == nil {
			continue
		}

		entry.complete = true
		s.enqueueIfReady(blockNumber, entry)
		delete(s.pendingBlocks, blockNumber)
	}
}

// recordQueueState must be called with s.mu held; the metrics port is non-reentrant.
func (s *BlockEventsQueue) recordQueueState(ctx context.Context) {
	s.syncMetrics.RecordQueueState(ctx, int64(len(s.readyList)), int64(len(s.pendingBlocks)))
}
