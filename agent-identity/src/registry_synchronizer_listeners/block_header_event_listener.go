package registry_synchronizer_listeners

import (
	"fmt"

	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type BlockHeaderEventListener struct {
	logger *zap.Logger

	blockEventHandlers                  []interfaces.BlockEventHandlerInterface
	blockchainEventSubscriptionProvider interfaces.BlockchainEventSubscriptionProviderInterface

	eventChannel chan *types.Header
	subscription ethereum.Subscription
}

type newBlockHeaderEventListenerParams struct {
	dig.In

	Logger *zap.Logger

	BlockEventHandlers                  []interfaces.BlockEventHandlerInterface `group:"block_event_handlers"`
	BlockchainEventSubscriptionProvider interfaces.BlockchainEventSubscriptionProviderInterface
}

func NewBlockHeaderEventListener(
	params newBlockHeaderEventListenerParams,
) *BlockHeaderEventListener {
	return &BlockHeaderEventListener{
		logger: params.Logger,

		blockEventHandlers:                  params.BlockEventHandlers,
		blockchainEventSubscriptionProvider: params.BlockchainEventSubscriptionProvider,

		eventChannel: make(chan *types.Header),
		subscription: nil,
	}
}

func (listener *BlockHeaderEventListener) SubscribeAll() error {
	listener.logger.Debug("subscribing to ethereum block headers")

	headers := make(chan *types.Header)
	subscription, err := listener.blockchainEventSubscriptionProvider.StartNewHeadSubscription(headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}

	listener.logger.Info("listening for new ethereum block headers")
	listener.eventChannel = headers
	listener.subscription = subscription

	return nil
}

func (listener *BlockHeaderEventListener) Listen() error {
	if listener.subscription == nil {
		return fmt.Errorf("no active block header subscriptions")
	}

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	events := listener.eventChannel

	wg.Add(1)
	go func(subscription ethereum.Subscription, eventChan <-chan *types.Header) {
		defer wg.Done()

		for {
			select {
			case header, ok := <-eventChan:
				if !ok {
					return
				}
				if header == nil || header.Number == nil {
					continue
				}

				blockNumber := header.Number.Uint64()
				if blockNumber == 0 {
					continue
				}

				completedBlock := blockNumber - 1
				for _, handler := range listener.blockEventHandlers {
					if err := handler.HandleBlockProcessed(completedBlock); err != nil {
						listener.logger.Error(
							"failed to handle processed block",
							zap.Uint64("block_number", completedBlock),
							zap.Error(err),
						)
					}
				}
			case err := <-listener.subscription.Err():
				select {
				case errCh <- fmt.Errorf("block header subscription error: %w", err):
				default:
				}
				return
			}
		}
	}(listener.subscription, events)

	err := <-errCh
	listener.subscription.Unsubscribe()
	wg.Wait()

	return err
}
