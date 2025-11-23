package subscribers

import (
	"context"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentCollectionAgentSignerUpdatedSubscriberInterface interface {
	SubscribeAgentSignerUpdated(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1AgentSignerUpdated) (ethereum.Subscription, error)
}
