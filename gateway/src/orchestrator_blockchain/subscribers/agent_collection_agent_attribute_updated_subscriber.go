package subscribers

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AgentCollectionAgentAttributeUpdatedSubscriber struct {
	logger *zap.Logger
	client *ethclient.Client
}

type newAgentCollectionAgentAttributeUpdatedSubscriberParams struct {
	dig.In

	Logger *zap.Logger
	Client *ethclient.Client `name:"ws"`
}

func NewAgentCollectionAgentAttributeUpdatedSubscriber(
	params newAgentCollectionAgentAttributeUpdatedSubscriberParams,
) *AgentCollectionAgentAttributeUpdatedSubscriber {
	return &AgentCollectionAgentAttributeUpdatedSubscriber{
		logger: params.Logger,
		client: params.Client,
	}
}

func (s *AgentCollectionAgentAttributeUpdatedSubscriber) SubscribeAgentAttributeUpdated(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1AgentAttributeUpdated) (ethereum.Subscription, error) {
	agentCollection, err := abi.NewAgentCollectionV1(agentCollectionAddress, s.client)
	if err != nil {
		return nil, err
	}

	sub, err := agentCollection.WatchAgentAttributeUpdated(nil, logs, nil, nil)
	if err != nil {
		return nil, err
	}

	s.logger.Info(
		"Subscribed to AgentCollection AgentAttributeUpdated events",
		zap.String("address", agentCollectionAddress.Hex()),
	)

	return sub, nil
}
