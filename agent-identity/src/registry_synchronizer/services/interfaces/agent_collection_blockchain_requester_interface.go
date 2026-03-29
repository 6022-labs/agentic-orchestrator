package interfaces

import (
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionBlockchainRequesterInterface interface {
	GetAgentOwner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error)
	GetAgentWalletsOfProposal(collectionAddress common.Address, proposalId *big.Int) ([]common.Address, error)
	GetAgentPreviousSigners(collectionAddress common.Address, collectionTokenId big.Int) ([]common.Address, error)
	GetAgentInformation(collectionAddress common.Address, collectionTokenId big.Int) (*abi.AgentInformation, error)

	GetAgentCollectionV1MintedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1Minted, error)
	GetAgentCollectionV1TransferEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1Transfer, error)
	GetAgentCollectionV1AgentImageUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1AgentImageUpdated, error)
	GetAgentCollectionV1AgentSignerUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1AgentSignerUpdated, error)
	GetAgentCollectionV1AgentAttributeUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1AgentAttributeUpdated, error)

	GetAgentCollectionV1MintProposalCreatedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1MintProposalCreated, error)
	GetAgentCollectionV1MintProposalRemovedEvent(from, to uint64, fromAddresses []common.Address) ([]*abi.AgentCollectionV1MintProposalRemoved, error)
}
