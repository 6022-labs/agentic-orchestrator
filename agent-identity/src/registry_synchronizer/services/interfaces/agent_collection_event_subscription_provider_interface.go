package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionEventSubscriptionProviderInterface interface {
	StartAgentAttributeUpdatedSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1AgentAttributeUpdated) (ethereum.Subscription, error)
	StartAgentSignerUpdatedSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1AgentSignerUpdated) (ethereum.Subscription, error)
	StartAgentImageUpdatedSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1AgentImageUpdated) (ethereum.Subscription, error)

	StartMintProposalCreatedSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1MintProposalCreated) (ethereum.Subscription, error)
	StartMintProposalRemovedSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1MintProposalRemoved) (ethereum.Subscription, error)
	StartMintedSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1Minted) (ethereum.Subscription, error)

	StartTransferSubscription(addresses []common.Address, logs chan<- *abi.AgentCollectionV1Transfer) (ethereum.Subscription, error)
}
