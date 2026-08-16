package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionsManagerEventSubscriptionProviderInterface interface {
	StartCollectionCreatedSubscription(
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionsManagerCollectionCreated, ethereum.Subscription, error)
}
