package subscribers

import (
	"context"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionMintedSubscriberInterface interface {
	SubscribeMinted(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1Minted) (ethereum.Subscription, error)
}
