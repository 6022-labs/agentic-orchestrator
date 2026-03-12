package interfaces

import "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"

type QueueItem struct {
	BlockNumber                                  uint64
	AgentCollectionV1MintedEvents                []*abi.AgentCollectionV1Minted
	AgentCollectionV1TransferEvents              []*abi.AgentCollectionV1Transfer
	AgentCollectionV1AgentImageUpdatedEvents     []*abi.AgentCollectionV1AgentImageUpdated
	AgentCollectionV1AgentSignerUpdatedEvents    []*abi.AgentCollectionV1AgentSignerUpdated
	AgentCollectionV1MintProposalCreatedEvents   []*abi.AgentCollectionV1MintProposalCreated
	AgentCollectionV1MintProposalRemovedEvents   []*abi.AgentCollectionV1MintProposalRemoved
	AgentCollectionV1AgentAttributeUpdatedEvents []*abi.AgentCollectionV1AgentAttributeUpdated

	AgentCollectionsManagerCollectionCreatedEvents []*abi.AgentCollectionsManagerCollectionCreated
}

type BlockEventsQueueInterface interface {
	BlockEventHandlerInterface

	AddAgentCollectionV1MintedEvents(blockNumber uint64, events []*abi.AgentCollectionV1Minted) error
	AddAgentCollectionV1TransferEvents(blockNumber uint64, events []*abi.AgentCollectionV1Transfer) error
	AddAgentCollectionV1AgentImageUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentImageUpdated) error
	AddAgentCollectionV1AgentSignerUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentSignerUpdated) error
	AddAgentCollectionV1MintProposalCreatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1MintProposalCreated) error
	AddAgentCollectionV1MintProposalRemovedEvents(blockNumber uint64, events []*abi.AgentCollectionV1MintProposalRemoved) error
	AddAgentCollectionV1AgentAttributeUpdatedEvents(blockNumber uint64, events []*abi.AgentCollectionV1AgentAttributeUpdated) error

	AddAgentCollectionsManagerCollectionCreatedEvents(blockNumber uint64, events []*abi.AgentCollectionsManagerCollectionCreated) error

	// Dequeue returns the next block to dequeue, ordered by block number ascending.
	Dequeue() (uint64, *QueueItem)
}
