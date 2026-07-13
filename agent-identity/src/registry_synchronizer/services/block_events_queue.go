package services

import (
	"sort"
	"sync"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/settings"
	"go.uber.org/zap"
)

const defaultBlockFinalizationDelay = 2 * time.Second

type BlockEventsQueue struct {
	logger       *zap.Logger
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
	syncSettings *settings.SyncSettings,
) *BlockEventsQueue {
	return &BlockEventsQueue{
		logger:       logger,
		syncSettings: syncSettings,

		entries:       make(map[uint64]*queueEntry),
		pendingBlocks: make(map[uint64]*queueEntry),
	}
}

func (s *BlockEventsQueue) AddAgentCollectionV1MintedEvents(blockNumber uint64, events []*abi.AgentCollectionV1Minted) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintedEvents = append(entry.data.AgentCollectionV1MintedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1TransferEvents(blockNumber uint64, events []*abi.AgentCollectionV1Transfer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1TransferEvents = append(entry.data.AgentCollectionV1TransferEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentImageUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentImageUpdated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentImageUpdatedEvents = append(entry.data.AgentCollectionV1AgentImageUpdatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1MintProposalCreatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1MintProposalCreated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintProposalCreatedEvents = append(entry.data.AgentCollectionV1MintProposalCreatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1MintProposalRemovedEvents(blockNumber uint64, events []*abi.AgentCollectionV1MintProposalRemoved) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintProposalRemovedEvents = append(entry.data.AgentCollectionV1MintProposalRemovedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentAttributeUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentAttributeUpdated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentAttributeUpdatedEvents = append(entry.data.AgentCollectionV1AgentAttributeUpdatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentWalletAddedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentWalletAdded) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentWalletAddedEvents = append(entry.data.AgentCollectionV1AgentWalletAddedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionV1AgentWalletRemovedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentWalletRemoved) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentWalletRemovedEvents = append(entry.data.AgentCollectionV1AgentWalletRemovedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) AddAgentCollectionsManagerCollectionCreatedEvents(blockNumber uint64, events []*abi.AgentCollectionsManagerCollectionCreated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionsManagerCollectionCreatedEvents = append(entry.data.AgentCollectionsManagerCollectionCreatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *BlockEventsQueue) Dequeue() (uint64, *interfaces.QueueItem) {
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

	s.readyList = append(s.readyList, blockNumber)
	entry.enqueued = true
	delete(s.pendingBlocks, blockNumber)
}

func (s *BlockEventsQueue) trackPendingBlock(blockNumber uint64, entry *queueEntry) {
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
	s.scheduleBlockFinalization(blockNumber)
}

func (s *BlockEventsQueue) scheduleBlockFinalization(blockNumber uint64) {
	if blockNumber > s.latestBlockSeen {
		s.latestBlockSeen = blockNumber
	}

	if s.timer != nil {
		return
	}

	s.timer = time.AfterFunc(s.getDeliveryBuffer(), func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.timer = nil

		completedBlock := s.latestBlockSeen
		if completedBlock > s.latestCompletedBlock {
			s.latestCompletedBlock = completedBlock
		}

		s.completePendingBlocks()
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
