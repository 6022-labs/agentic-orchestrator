package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionEventSubscriptionProviderInterface interface {
	StartAgentAttributeUpdatedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1AgentAttributeUpdated, ethereum.Subscription, error)
	StartAgentAddressAddedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1AgentAddressAdded, ethereum.Subscription, error)
	StartAgentAddressRemovedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1AgentAddressRemoved, ethereum.Subscription, error)
	StartAgentImageUpdatedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1AgentImageUpdated, ethereum.Subscription, error)
	StartAgentImageRemovedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1AgentImageRemoved, ethereum.Subscription, error)

	StartMintProposalCreatedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1MintProposalCreated, ethereum.Subscription, error)
	StartMintProposalRemovedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1MintProposalRemoved, ethereum.Subscription, error)
	StartMintedSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1Minted, ethereum.Subscription, error)

	StartTransferSubscription(
		ctx context.Context,
		addresses []common.Address,
	) (<-chan *abi.AgentCollectionV1Transfer, ethereum.Subscription, error)
}
