package interfaces

import (
	"math/big"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionRequesterInterface interface {
	GetName(collectionAddress common.Address) (*string, error)
	GetSymbol(collectionAddress common.Address) (*string, error)
	GetNextTokenId(collectionAddress common.Address) (*big.Int, error)
	GetAgentOwner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error)
	GetAgentSigner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error)
	GetAgentWallets(collectionAddress common.Address, collectionTokenId big.Int) ([]common.Address, error)
	GetAgentWalletsOfProposal(collectionAddress common.Address, proposalId *big.Int) ([]common.Address, error)
	GetAgentPreviousSigners(collectionAddress common.Address, collectionTokenId big.Int) ([]common.Address, error)
	GetAgentInformation(collectionAddress common.Address, collectionTokenId big.Int) (*abi.AgentInformation, error)
	GetAddOrUpdateAgentAttributeInputsFromTransaction(collectionAddress common.Address, txHash common.Hash) (tokenId *big.Int, key *string, value *string, err error)
}
