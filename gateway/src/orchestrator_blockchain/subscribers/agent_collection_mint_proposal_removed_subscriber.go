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

type AgentCollectionMintProposalRemovedSubscriber struct {
	logger *zap.Logger
	client *ethclient.Client
}

type newAgentCollectionMintProposalRemovedSubscriberParams struct {
	dig.In

	Logger *zap.Logger
	Client *ethclient.Client `name:"ws"`
}

func NewAgentCollectionMintProposalRemovedSubscriber(
	params newAgentCollectionMintProposalRemovedSubscriberParams,
) *AgentCollectionMintProposalRemovedSubscriber {
	return &AgentCollectionMintProposalRemovedSubscriber{
		logger: params.Logger,
		client: params.Client,
	}
}

func (s *AgentCollectionMintProposalRemovedSubscriber) SubscribeMintProposalRemoved(ctx context.Context, agentCollectionAddress common.Address, logs chan<- *abi.AgentCollectionV1MintProposalRemoved) (ethereum.Subscription, error) {
	agentCollection, err := abi.NewAgentCollectionV1(agentCollectionAddress, s.client)
	if err != nil {
		return nil, err
	}

	sub, err := agentCollection.WatchMintProposalRemoved(nil, logs, nil, nil)
	if err != nil {
		return nil, err
	}

	s.logger.Info(
		"Subscribed to AgentCollection MintProposalRemoved events",
		zap.String("address", agentCollectionAddress.Hex()),
	)

	return sub, nil
}
