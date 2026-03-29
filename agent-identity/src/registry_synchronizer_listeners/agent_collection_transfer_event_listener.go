package registry_synchronizer_listeners

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionTransferListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	agentCollectionRepository                repositories.AgentCollectionRepositoryInterface
	transferEventHandler                     interfaces.TransferEventHandlerInterface
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface

	trackedAddresses []common.Address
}

func NewAgentCollectionTransferListener(
	logger *zap.Logger,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	transferEventHandler interfaces.TransferEventHandlerInterface,
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface,
) *AgentCollectionTransferListener {
	return &AgentCollectionTransferListener{
		AbstractUpdatableSubscriptionListener: NewAbstractUpdatableSubscriptionListener(),

		logger: logger,

		agentCollectionRepository:                agentCollectionRepository,
		transferEventHandler:                     transferEventHandler,
		agentCollectionEventSubscriptionProvider: agentCollectionEventSubscriptionProvider,

		trackedAddresses: []common.Address{},
	}
}

func (listener *AgentCollectionTransferListener) Listen() error {
	return listener.AbstractUpdatableSubscriptionListener.Listen()
}

func (listener *AgentCollectionTransferListener) SubscribeAll() error {
	collections, err := listener.agentCollectionRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get agent collections: %w", err)
	}

	collectionAddresses := make([]common.Address, 0, len(collections))
	for _, collection := range collections {
		collectionAddresses = append(collectionAddresses, common.HexToAddress(collection.Address))
	}

	if len(collectionAddresses) == 0 {
		listener.logger.Info("no agent collections found, skipping subscription to AgentCollection.Transfer events")
		return nil
	}

	listener.mu.Lock()
	listener.trackedAddresses = append([]common.Address(nil), collectionAddresses...)
	addressesToSubscribe := append([]common.Address(nil), collectionAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartTransferSubscription(addressesToSubscribe)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.Transfer events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info("subscribed to AgentCollection.Transfer events", zap.Int("collectionCount", len(addressesToSubscribe)))

	return nil
}

func (listener *AgentCollectionTransferListener) Subscribe(collectionAddress common.Address) error {
	listener.mu.Lock()
	for _, trackedAddress := range listener.trackedAddresses {
		if trackedAddress == collectionAddress {
			listener.mu.Unlock()
			return nil
		}
	}

	listener.trackedAddresses = append(listener.trackedAddresses, collectionAddress)
	addressesToSubscribe := append([]common.Address(nil), listener.trackedAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartTransferSubscription(addressesToSubscribe)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.Transfer events for new address: %w", err)
	}

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.oldSubscriptions = append(listener.oldSubscriptions, listener.subscription)
	}
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info("started new AgentCollection.Transfer subscription for collection", zap.String("collectionAddress", collectionAddress.Hex()))

	return nil
}

func (listener *AgentCollectionTransferListener) startSubscriptionHandler(subscription ethereum.Subscription, eventChan <-chan *abi.AgentCollectionV1Transfer) {
	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()

		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					return
				}
				if event == nil {
					continue
				}

				if err := listener.transferEventHandler.Handle(event); err != nil {
					listener.logger.Error("failed to handle AgentCollection.Transfer event", zap.Any("event", event), zap.Error(err))
				}

				listener.mu.Lock()
				if event.Raw.BlockNumber > listener.lastBlockSeen {
					listener.lastBlockSeen = event.Raw.BlockNumber
				}
				listener.mu.Unlock()
			case err := <-subscription.Err():
				if err == nil {
					return
				}
				select {
				case listener.errCh <- fmt.Errorf("agent collection transfer subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
