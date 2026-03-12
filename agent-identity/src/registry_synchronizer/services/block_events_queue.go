package services

import (
	"sort"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
)

type UniswapV3BlockEventsQueue struct {
	mu                   sync.Mutex
	entries              map[uint64]*queueEntry
	readyList            []uint64
	pendingBlocks        map[uint64]*queueEntry
	latestCompletedBlock uint64
}

type queueEntry struct {
	data      *interfaces.QueueItem
	hasEvents bool
	complete  bool
	enqueued  bool
}

func NewUniswapV3BlockEventsQueue() *UniswapV3BlockEventsQueue {
	return &UniswapV3BlockEventsQueue{
		entries:       make(map[uint64]*queueEntry),
		pendingBlocks: make(map[uint64]*queueEntry),
	}
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1MintedEvents(blockNumber uint64, events []*abi.AgentCollectionV1Minted) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintedEvents = append(entry.data.AgentCollectionV1MintedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1TransferEvents(blockNumber uint64, events []*abi.AgentCollectionV1Transfer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1TransferEvents = append(entry.data.AgentCollectionV1TransferEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1AgentImageUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentImageUpdated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentImageUpdatedEvents = append(entry.data.AgentCollectionV1AgentImageUpdatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1AgentSignerUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentSignerUpdated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentSignerUpdatedEvents = append(entry.data.AgentCollectionV1AgentSignerUpdatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1MintProposalCreatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1MintProposalCreated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintProposalCreatedEvents = append(entry.data.AgentCollectionV1MintProposalCreatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1MintProposalRemovedEvents(blockNumber uint64, events []*abi.AgentCollectionV1MintProposalRemoved) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1MintProposalRemovedEvents = append(entry.data.AgentCollectionV1MintProposalRemovedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionV1AgentAttributeUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentAttributeUpdated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionV1AgentAttributeUpdatedEvents = append(entry.data.AgentCollectionV1AgentAttributeUpdatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) AddAgentCollectionsManagerCollectionCreatedEvents(blockNumber uint64, events []*abi.AgentCollectionsManagerCollectionCreated) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := s.ensureEntry(blockNumber)
	entry.data.AgentCollectionsManagerCollectionCreatedEvents = append(entry.data.AgentCollectionsManagerCollectionCreatedEvents, events...)
	s.trackPendingBlock(blockNumber, entry)

	return nil
}

func (s *UniswapV3BlockEventsQueue) MarkBlockAsComplete(blockNumber uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if blockNumber > s.latestCompletedBlock {
		s.latestCompletedBlock = blockNumber
	}

	s.completePendingBlocks()

	return nil
}

func (s *UniswapV3BlockEventsQueue) Dequeue() (uint64, *interfaces.QueueItem) {
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

func (s *UniswapV3BlockEventsQueue) ensureEntry(blockNumber uint64) *queueEntry {
	entry, exists := s.entries[blockNumber]
	if !exists {
		entry = &queueEntry{
			data: &interfaces.QueueItem{},
		}
		s.entries[blockNumber] = entry
	}

	return entry
}

func (s *UniswapV3BlockEventsQueue) enqueueIfReady(blockNumber uint64, entry *queueEntry) {
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

func (s *UniswapV3BlockEventsQueue) trackPendingBlock(blockNumber uint64, entry *queueEntry) {
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
}

func (s *UniswapV3BlockEventsQueue) completePendingBlocks() {
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
