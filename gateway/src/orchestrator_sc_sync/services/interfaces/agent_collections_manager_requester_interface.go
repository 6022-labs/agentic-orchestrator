package interfaces

import "github.com/ethereum/go-ethereum/common"

type AgentCollectionsManagerRequesterInterface interface {
	GetCollectionAddresses() ([]common.Address, error)
}
