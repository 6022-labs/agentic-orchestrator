package registry_synchronizer_listeners

import (
	"fmt"
	"slices"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionAgentSignerUpdatedListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	agentCollectionRepository                repositories.AgentCollectionRepositoryInterface
	agentSignerUpdatedEventHandler           interfaces.AgentSignerUpdatedEventHandlerInterface
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface

	trackedAddresses []common.Address
}

func NewAgentCollectionAgentSignerUpdatedListener(
	logger *zap.Logger,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	agentSignerUpdatedEventHandler interfaces.AgentSignerUpdatedEventHandlerInterface,
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface,
) *AgentCollectionAgentSignerUpdatedListener {
	return &AgentCollectionAgentSignerUpdatedListener{
		AbstractUpdatableSubscriptionListener: NewAbstractUpdatableSubscriptionListener(),

		logger: logger,

		agentCollectionRepository:                agentCollectionRepository,
		agentSignerUpdatedEventHandler:           agentSignerUpdatedEventHandler,
		agentCollectionEventSubscriptionProvider: agentCollectionEventSubscriptionProvider,

		trackedAddresses: []common.Address{},
	}
}

func (listener *AgentCollectionAgentSignerUpdatedListener) Listen() error {
	return listener.AbstractUpdatableSubscriptionListener.Listen()
}

func (listener *AgentCollectionAgentSignerUpdatedListener) SubscribeAll() error {
	collections, err := listener.agentCollectionRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get agent collections: %w", err)
	}

	collectionAddresses := make([]common.Address, 0, len(collections))
	for _, collection := range collections {
		collectionAddresses = append(collectionAddresses, common.HexToAddress(collection.Address))
	}

	if len(collectionAddresses) == 0 {
		listener.logger.Info("no agent collections found, skipping subscription to AgentCollection.AgentSignerUpdated events")
		return nil
	}

	listener.mu.Lock()
	listener.trackedAddresses = append([]common.Address(nil), collectionAddresses...)
	addressesToSubscribe := append([]common.Address(nil), collectionAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentSignerUpdatedSubscription(addressesToSubscribe)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.AgentSignerUpdated events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info("subscribed to AgentCollection.AgentSignerUpdated events", zap.Int("collectionCount", len(addressesToSubscribe)))

	return nil
}

func (listener *AgentCollectionAgentSignerUpdatedListener) Subscribe(collectionAddress common.Address) error {
	listener.mu.Lock()
	if slices.Contains(listener.trackedAddresses, collectionAddress) {
		listener.mu.Unlock()
		return nil
	}

	listener.trackedAddresses = append(listener.trackedAddresses, collectionAddress)
	addressesToSubscribe := append([]common.Address(nil), listener.trackedAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentSignerUpdatedSubscription(addressesToSubscribe)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.AgentSignerUpdated events for new address: %w", err)
	}

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.oldSubscriptions = append(listener.oldSubscriptions, listener.subscription)
	}
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info("started new AgentCollection.AgentSignerUpdated subscription for collection", zap.String("collectionAddress", collectionAddress.Hex()))

	return nil
}

func (listener *AgentCollectionAgentSignerUpdatedListener) startSubscriptionHandler(subscription ethereum.Subscription, eventChan <-chan *abi.AgentCollectionV1AgentSignerUpdated) {
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

				if err := listener.agentSignerUpdatedEventHandler.Handle(event); err != nil {
					listener.logger.Error("failed to handle AgentCollection.AgentSignerUpdated event", zap.Any("event", event), zap.Error(err))
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
				case listener.errCh <- fmt.Errorf("agent collection agent signer updated subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
