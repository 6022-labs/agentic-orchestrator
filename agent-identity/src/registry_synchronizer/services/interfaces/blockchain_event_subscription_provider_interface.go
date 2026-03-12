package interfaces

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockchainEventSubscriptionProviderInterface interface {
	StartNewHeadSubscription(headers chan<- *types.Header) (ethereum.Subscription, error)
}
