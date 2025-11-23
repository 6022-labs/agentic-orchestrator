package subscribers

import (
	"context"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum"
)

type AgentCollectionsManagerCollectionCreatedSubscriberInterface interface {
	SubscribeCollectionCreated(ctx context.Context, logs chan<- *abi.AgentCollectionsManagerCollectionCreated) (ethereum.Subscription, error)
}
