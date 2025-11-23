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

type AgentCollectionAgentImageUpdatedSubscriber struct {
	logger *zap.Logger
	client *ethclient.Client
}

type newAgentCollectionAgentImageUpdatedSubscriberParams struct {
	dig.In

	Logger *zap.Logger
	Client *ethclient.Client `name:"ws"`
}

func NewAgentCollectionAgentImageUpdatedSubscriber(
	params newAgentCollectionAgentImageUpdatedSubscriberParams,
) *AgentCollectionAgentImageUpdatedSubscriber {
	return &AgentCollectionAgentImageUpdatedSubscriber{
		logger: params.Logger,
		client: params.Client,
	}
}

func (s *AgentCollectionAgentImageUpdatedSubscriber) SubscribeAgentImageUpdated(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1AgentImageUpdated) (ethereum.Subscription, error) {
	agentCollection, err := abi.NewAgentCollectionV1(agentCollectionAddress, s.client)
	if err != nil {
		return nil, err
	}

	sub, err := agentCollection.WatchAgentImageUpdated(nil, logs, nil, nil)
	if err != nil {
		return nil, err
	}

	s.logger.Info(
		"Subscribed to AgentCollection AgentImageUpdated events",
		zap.String("address", agentCollectionAddress.Hex()),
	)

	return sub, nil
}
