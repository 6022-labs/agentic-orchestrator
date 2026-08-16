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

type AgentCollectionAgentAddressAddedListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	agentCollectionRepository                repositories.AgentCollectionRepositoryInterface
	blockEventsQueue                         interfaces.BlockEventsQueueInterface
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface

	trackedAddresses []common.Address
}

func NewAgentCollectionAgentAddressAddedListener(
	logger *zap.Logger,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface,
) *AgentCollectionAgentAddressAddedListener {
	el := &AgentCollectionAgentAddressAddedListener{
		logger: logger,

		agentCollectionRepository:                agentCollectionRepository,
		blockEventsQueue:                         blockEventsQueue,
		agentCollectionEventSubscriptionProvider: agentCollectionEventSubscriptionProvider,

		trackedAddresses: []common.Address{},
	}
	el.AbstractUpdatableSubscriptionListener = NewAbstractUpdatableSubscriptionListener(el.rebuildSubscription)

	return el
}

func (listener *AgentCollectionAgentAddressAddedListener) Listen() error {
	return listener.AbstractUpdatableSubscriptionListener.Listen()
}

func (listener *AgentCollectionAgentAddressAddedListener) SubscribeAll() error {
	collections, err := listener.agentCollectionRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get agent collections: %w", err)
	}

	collectionAddresses := make([]common.Address, 0, len(collections))
	for _, collection := range collections {
		collectionAddresses = append(collectionAddresses, common.HexToAddress(collection.Address))
	}

	if len(collectionAddresses) == 0 {
		listener.logger.Info(
			"no agent collections found, skipping subscription to AgentCollection.AgentAddressAdded events",
		)
		return nil
	}

	listener.mu.Lock()
	listener.trackedAddresses = append([]common.Address(nil), collectionAddresses...)
	addressesToSubscribe := append([]common.Address(nil), collectionAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentAddressAddedSubscription(
		addressesToSubscribe,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.AgentAddressAdded events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info(
		"subscribed to AgentCollection.AgentAddressAdded events",
		zap.Int("collectionCount", len(addressesToSubscribe)),
	)

	return nil
}

func (listener *AgentCollectionAgentAddressAddedListener) Subscribe(collectionAddress common.Address) error {
	listener.mu.Lock()
	if slices.Contains(listener.trackedAddresses, collectionAddress) {
		listener.mu.Unlock()
		return nil
	}
	listener.trackedAddresses = append(listener.trackedAddresses, collectionAddress)
	listener.mu.Unlock()

	listener.markNeedsRebuild()
	return nil
}

func (listener *AgentCollectionAgentAddressAddedListener) rebuildSubscription() error {
	listener.mu.Lock()
	addressesToSubscribe := append([]common.Address(nil), listener.trackedAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentAddressAddedSubscription(
		addressesToSubscribe,
	)
	if err != nil {
		return fmt.Errorf("failed to rebuild AgentCollection.AgentAddressAdded subscription: %w", err)
	}

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.oldSubscriptions = append(listener.oldSubscriptions, listener.subscription)
	}
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info(
		"rebuilt AgentCollection.AgentAddressAdded subscription",
		zap.Int("collectionCount", len(addressesToSubscribe)),
	)
	return nil
}

func (listener *AgentCollectionAgentAddressAddedListener) startSubscriptionHandler(
	subscription ethereum.Subscription,
	eventChan <-chan *abi.AgentCollectionV1AgentAddressAdded,
) {
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

				if err := listener.blockEventsQueue.AddAgentCollectionV1AgentAddressAddedEvents(
					event.Raw.BlockNumber,
					[]*abi.AgentCollectionV1AgentAddressAdded{event},
				); err != nil {
					listener.logger.Error(
						"failed to enqueue AgentCollection.AgentAddressAdded event",
						zap.Uint64("blockNumber", event.Raw.BlockNumber),
						zap.Error(err),
					)
				}

				listener.onBlockSeen(event.Raw.BlockNumber)
			case err := <-subscription.Err():
				if err == nil {
					return
				}
				select {
				case listener.errCh <- fmt.Errorf("agent collection agent address added subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
