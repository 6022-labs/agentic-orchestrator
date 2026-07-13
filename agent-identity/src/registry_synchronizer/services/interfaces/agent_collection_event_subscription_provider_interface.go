package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionEventSubscriptionProviderInterface interface {
	StartAgentAttributeUpdatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentAttributeUpdated, ethereum.Subscription, error)
	StartAgentWalletAddedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentWalletAdded, ethereum.Subscription, error)
	StartAgentWalletRemovedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentWalletRemoved, ethereum.Subscription, error)
	StartAgentImageUpdatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentImageUpdated, ethereum.Subscription, error)

	StartMintProposalCreatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1MintProposalCreated, ethereum.Subscription, error)
	StartMintProposalRemovedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1MintProposalRemoved, ethereum.Subscription, error)
	StartMintedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1Minted, ethereum.Subscription, error)

	StartTransferSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1Transfer, ethereum.Subscription, error)
}
