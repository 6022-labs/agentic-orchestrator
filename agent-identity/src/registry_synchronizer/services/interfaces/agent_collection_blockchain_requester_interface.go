package interfaces

import (
	"context"

	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionBlockchainRequesterInterface interface {
	GetAgentOwner(
		ctx context.Context,
		collectionAddress common.Address,
		collectionTokenId big.Int,
	) (*common.Address, error)
	GetMintProposal(
		ctx context.Context,
		collectionAddress common.Address,
		proposalId *big.Int,
		atBlock uint64,
	) (*abi.MintProposal, error)
	GetAgentInformation(
		ctx context.Context,
		collectionAddress common.Address,
		collectionTokenId big.Int,
	) (*abi.AgentInformation, error)

	GetAgentCollectionV1MintedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1Minted, error)
	GetAgentCollectionV1TransferEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1Transfer, error)
	GetAgentCollectionV1AgentImageUpdatedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1AgentImageUpdated, error)
	GetAgentCollectionV1AgentImageRemovedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1AgentImageRemoved, error)
	GetAgentCollectionV1AgentAttributeUpdatedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1AgentAttributeUpdated, error)
	GetAgentCollectionV1AgentAddressAddedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1AgentAddressAdded, error)
	GetAgentCollectionV1AgentAddressRemovedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1AgentAddressRemoved, error)

	GetAgentCollectionV1MintProposalCreatedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1MintProposalCreated, error)
	GetAgentCollectionV1MintProposalRemovedEvent(
		ctx context.Context,
		from, to uint64,
		fromAddresses []common.Address,
	) ([]*abi.AgentCollectionV1MintProposalRemoved, error)
}
