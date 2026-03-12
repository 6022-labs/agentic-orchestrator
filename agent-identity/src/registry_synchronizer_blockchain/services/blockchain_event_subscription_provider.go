package services

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type BlockchainEventSubscriptionProvider struct {
	logger *zap.Logger
	client *ethclient.Client
}

type newBlockchainEventSubscriptionProviderParams struct {
	dig.In

	Logger *zap.Logger
	Client *ethclient.Client `name:"ws"`
}

func NewBlockchainEventSubscriptionProvider(
	params newBlockchainEventSubscriptionProviderParams,
) *BlockchainEventSubscriptionProvider {
	return &BlockchainEventSubscriptionProvider{
		logger: params.Logger,
		client: params.Client,
	}
}

func (s *BlockchainEventSubscriptionProvider) StartNewHeadSubscription(headers chan<- *types.Header) (ethereum.Subscription, error) {
	sub, err := s.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return nil, err
	}

	wrappedSub := event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	})

	s.logger.Info("Subscribed to new Ethereum block headers")

	return wrappedSub, nil
}
