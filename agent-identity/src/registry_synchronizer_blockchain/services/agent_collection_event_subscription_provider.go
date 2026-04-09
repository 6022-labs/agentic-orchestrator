package services

import (
	"context"
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AgentCollectionEventSubscriptionProvider struct {
	logger *zap.Logger
	client *ethclient.Client
}

type newAgentCollectionEventSubscriptionProviderParams struct {
	dig.In

	Logger *zap.Logger
	Client *ethclient.Client `name:"ws"`
}

func NewAgentCollectionEventSubscriptionProvider(
	params newAgentCollectionEventSubscriptionProviderParams,
) *AgentCollectionEventSubscriptionProvider {
	return &AgentCollectionEventSubscriptionProvider{
		logger: params.Logger,
		client: params.Client,
	}
}

func (s *AgentCollectionEventSubscriptionProvider) StartAgentAttributeUpdatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentAttributeUpdated, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["AgentAttributeUpdated"]
	if !ok {
		return nil, nil, fmt.Errorf("event AgentAttributeUpdated not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1AgentAttributeUpdated, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseAgentAttributeUpdated(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse AgentAttributeUpdated log: %w", parseErr)
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
		"Subscribed to AgentCollection AgentAttributeUpdated events",
	)

	return logs, wrappedSub, nil
}

func (s *AgentCollectionEventSubscriptionProvider) StartAgentSignerUpdatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentSignerUpdated, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["AgentSignerUpdated"]
	if !ok {
		return nil, nil, fmt.Errorf("event AgentSignerUpdated not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1AgentSignerUpdated, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseAgentSignerUpdated(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse AgentSignerUpdated log: %w", parseErr)
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
		"Subscribed to AgentCollection AgentSignerUpdated events",
	)

	return logs, wrappedSub, nil
}

func (s *AgentCollectionEventSubscriptionProvider) StartAgentImageUpdatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1AgentImageUpdated, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["AgentImageUpdated"]
	if !ok {
		return nil, nil, fmt.Errorf("event AgentImageUpdated not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1AgentImageUpdated, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseAgentImageUpdated(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse AgentImageUpdated log: %w", parseErr)
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
		"Subscribed to AgentCollection AgentImageUpdated events",
	)

	return logs, wrappedSub, nil
}

func (s *AgentCollectionEventSubscriptionProvider) StartMintProposalCreatedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1MintProposalCreated, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["MintProposalCreated"]
	if !ok {
		return nil, nil, fmt.Errorf("event MintProposalCreated not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1MintProposalCreated, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseMintProposalCreated(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse MintProposalCreated log: %w", parseErr)
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
		"Subscribed to AgentCollection MintProposalCreated events",
	)

	return logs, wrappedSub, nil
}

func (s *AgentCollectionEventSubscriptionProvider) StartMintProposalRemovedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1MintProposalRemoved, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["MintProposalRemoved"]
	if !ok {
		return nil, nil, fmt.Errorf("event MintProposalRemoved not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1MintProposalRemoved, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseMintProposalRemoved(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse MintProposalRemoved log: %w", parseErr)
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
		"Subscribed to AgentCollection MintProposalRemoved events",
	)

	return logs, wrappedSub, nil
}

func (s *AgentCollectionEventSubscriptionProvider) StartMintedSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1Minted, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["Minted"]
	if !ok {
		return nil, nil, fmt.Errorf("event Minted not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1Minted, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseMinted(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse Minted log: %w", parseErr)
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
		"Subscribed to AgentCollection Minted events",
	)

	return logs, wrappedSub, nil
}

func (s *AgentCollectionEventSubscriptionProvider) StartTransferSubscription(addresses []common.Address) (<-chan *abi.AgentCollectionV1Transfer, ethereum.Subscription, error) {
	abiSpec, err := abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	eventDef, ok := abiSpec.Events["Transfer"]
	if !ok {
		return nil, nil, fmt.Errorf("event Transfer not found in AgentCollectionV1 ABI")
	}

	query := ethereum.FilterQuery{Topics: [][]common.Hash{{eventDef.ID}}, Addresses: addresses}
	subscriptionLogs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, subscriptionLogs)
	if err != nil {
		return nil, nil, err
	}

	filterer, err := abi.NewAgentCollectionV1Filterer(common.Address{}, s.client)
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan *abi.AgentCollectionV1Transfer, 64)

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer close(logs)
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-subscriptionLogs:
				parsedEvent, parseErr := filterer.ParseTransfer(log)
				if parseErr != nil {
					return fmt.Errorf("failed to parse Transfer log: %w", parseErr)
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
		"Subscribed to AgentCollection Transfer events",
	)

	return logs, wrappedSub, nil
}
