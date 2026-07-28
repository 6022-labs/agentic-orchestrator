package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AgentCollectionsManagerEventSubscriptionProvider struct {
	logger                          *zap.Logger
	client                          *ethclient.Client
	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

type NewAgentCollectionsManagerEventSubscriptionProviderParams struct {
	dig.In

	Logger                          *zap.Logger
	Client                          *ethclient.Client `name:"ws"`
	AgentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings
}

func NewAgentCollectionsManagerEventSubscriptionProvider(
	params NewAgentCollectionsManagerEventSubscriptionProviderParams,
) *AgentCollectionsManagerEventSubscriptionProvider {
	return &AgentCollectionsManagerEventSubscriptionProvider{
		logger:                          params.Logger,
		client:                          params.Client,
		agentCollectionsManagerSettings: params.AgentCollectionsManagerSettings,
	}
}

func (s *AgentCollectionsManagerEventSubscriptionProvider) StartCollectionCreatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionsManagerCollectionCreated, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionsManagerMetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["CollectionCreated"]
	if !ok {
		return nil, nil, fmt.Errorf("event CollectionCreated not found in AgentCollectionsManager ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionsManagerFilterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionsManagerCollectionCreated, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseCollectionCreated(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse CollectionCreated log: %w", parseErr)
				}
				select {
				case logs <- parsedEvent:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	})

	s.logger.Debug(
		"Subscribed to AgentCollectionsManager CollectionCreated events",
	)

	return logs, wrappedSub, nil
}
