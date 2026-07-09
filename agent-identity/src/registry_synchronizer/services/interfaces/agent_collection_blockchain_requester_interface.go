package interfaces

import (
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionBlockchainRequesterInterface interface {
	GetAgentOwner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error)
	GetMintProposal(collectionAddress common.Address, proposalId *big.Int, atBlock uint64) (*abi.MintProposal, error)
	GetAgentInformation(collectionAddress common.Address, collectionTokenId big.Int) (*abi.AgentInformation, error)

	GetAgentCollectionV1MintedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1Minted, error)
	GetAgentCollectionV1TransferEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1Transfer, error)
	GetAgentCollectionV1AgentImageUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1AgentImageUpdated, error)
	GetAgentCollectionV1AgentAttributeUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1AgentAttributeUpdated, error)

	GetAgentCollectionV1MintProposalCreatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1MintProposalCreated, error)
	GetAgentCollectionV1MintProposalRemovedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1MintProposalRemoved, error)
}
