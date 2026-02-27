package subscribers

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionMintProposalCreatedSubscriberInterface interface {
	SubscribeMintProposalCreated(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1MintProposalCreated) (ethereum.Subscription, error)
}
