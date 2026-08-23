package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
)

type QueueItem struct {
	BlockNumber                                  uint64
	AgentCollectionV1MintedEvents                []*abi.AgentCollectionV1Minted
	AgentCollectionV1TransferEvents              []*abi.AgentCollectionV1Transfer
	AgentCollectionV1AgentAddressAddedEvents     []*abi.AgentCollectionV1AgentAddressAdded
	AgentCollectionV1AgentImageUpdatedEvents     []*abi.AgentCollectionV1AgentImageUpdated
	AgentCollectionV1AgentImageRemovedEvents     []*abi.AgentCollectionV1AgentImageRemoved
	AgentCollectionV1AgentAddressRemovedEvents   []*abi.AgentCollectionV1AgentAddressRemoved
	AgentCollectionV1MintProposalCreatedEvents   []*abi.AgentCollectionV1MintProposalCreated
	AgentCollectionV1MintProposalRemovedEvents   []*abi.AgentCollectionV1MintProposalRemoved
	AgentCollectionV1AgentAttributeUpdatedEvents []*abi.AgentCollectionV1AgentAttributeUpdated

	AgentCollectionsManagerCollectionCreatedEvents []*abi.AgentCollectionsManagerCollectionCreated
}

type BlockEventsQueueInterface interface {
	AddAgentCollectionV1MintedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1Minted,
	) error
	AddAgentCollectionV1TransferEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1Transfer,
	) error
	AddAgentCollectionV1AgentAddressAddedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1AgentAddressAdded,
	) error
	AddAgentCollectionV1AgentImageUpdatedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1AgentImageUpdated,
	) error
	AddAgentCollectionV1AgentImageRemovedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1AgentImageRemoved,
	) error
	AddAgentCollectionV1AgentAddressRemovedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1AgentAddressRemoved,
	) error
	AddAgentCollectionV1MintProposalCreatedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1MintProposalCreated,
	) error
	AddAgentCollectionV1MintProposalRemovedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1MintProposalRemoved,
	) error
	AddAgentCollectionV1AgentAttributeUpdatedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionV1AgentAttributeUpdated,
	) error

	AddAgentCollectionsManagerCollectionCreatedEvents(
		ctx context.Context,
		blockNumber uint64,
		events []*abi.AgentCollectionsManagerCollectionCreated,
	) error

	// Dequeue returns the next block to dequeue, ordered by block number ascending.
	Dequeue(ctx context.Context) (uint64, *QueueItem)
}
