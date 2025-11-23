package subscribers

import (
	"context"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AgentCollectionAgentSignerUpdatedSubscriber struct {
	logger *zap.Logger
	client *ethclient.Client
}

type newAgentCollectionAgentSignerUpdatedSubscriberParams struct {
	dig.In

	Logger *zap.Logger
	Client *ethclient.Client `name:"ws"`
}

func NewAgentCollectionAgentSignerUpdatedSubscriber(
	params newAgentCollectionAgentSignerUpdatedSubscriberParams,
) *AgentCollectionAgentSignerUpdatedSubscriber {
	return &AgentCollectionAgentSignerUpdatedSubscriber{
		logger: params.Logger,
		client: params.Client,
	}
}

func (s *AgentCollectionAgentSignerUpdatedSubscriber) SubscribeAgentSignerUpdated(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1AgentSignerUpdated) (ethereum.Subscription, error) {
	agentCollection, err := abi.NewAgentCollectionV1(agentCollectionAddress, s.client)
	if err != nil {
		return nil, err
	}

	sub, err := agentCollection.WatchAgentSignerUpdated(nil, logs, nil)
	if err != nil {
		return nil, err
	}

	s.logger.Info(
		"Subscribed to AgentCollection AgentSignerUpdated events",
		zap.String("address", agentCollectionAddress.Hex()),
	)

	return sub, nil
}
