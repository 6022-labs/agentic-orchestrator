package registry_synchronizer_listeners

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AgentCollectionsManagerCollectionCreatedListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings

	collectionListeners []interfaces.CollectionEventSubscriberInterface

	collectionCreatedEventHandler                    interfaces.CollectionCreatedEventHandlerInterface
	agentCollectionsManagerEventSubscriptionProvider interfaces.AgentCollectionsManagerEventSubscriptionProviderInterface
}

type newAgentCollectionsManagerCollectionCreatedListenerParams struct {
	dig.In

	Logger                                           *zap.Logger
	AgentCollectionsManagerSettings                  *settings.AgentCollectionsManagerSettings
	CollectionListeners                              []interfaces.CollectionEventSubscriberInterface `group:"collection_event_subscribers"`
	CollectionCreatedEventHandler                    interfaces.CollectionCreatedEventHandlerInterface
	AgentCollectionsManagerEventSubscriptionProvider interfaces.AgentCollectionsManagerEventSubscriptionProviderInterface
}

func NewAgentCollectionsManagerCollectionCreatedListener(
	params newAgentCollectionsManagerCollectionCreatedListenerParams,
) *AgentCollectionsManagerCollectionCreatedListener {
	return &AgentCollectionsManagerCollectionCreatedListener{
		AbstractUpdatableSubscriptionListener: NewAbstractUpdatableSubscriptionListener(),

		logger: params.Logger,

		agentCollectionsManagerSettings: params.AgentCollectionsManagerSettings,

		collectionListeners: params.CollectionListeners,

		collectionCreatedEventHandler:                    params.CollectionCreatedEventHandler,
		agentCollectionsManagerEventSubscriptionProvider: params.AgentCollectionsManagerEventSubscriptionProvider,
	}
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) Listen() error {
	return listener.AbstractUpdatableSubscriptionListener.Listen()
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) SubscribeAll() error {
	events := make(chan *abi.AgentCollectionsManagerCollectionCreated, 64)

	subscription, err := listener.agentCollectionsManagerEventSubscriptionProvider.StartCollectionCreatedSubscription(
		[]common.Address{listener.agentCollectionsManagerSettings.SmartContractAddress},
		events,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollectionsManager.CollectionCreated events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)

	listener.logger.Info(
		"subscribed to AgentCollectionsManager.CollectionCreated events",
		zap.String("agentCollectionsManagerAddress", listener.agentCollectionsManagerSettings.SmartContractAddress.Hex()),
	)

	return nil
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) startSubscriptionHandler(subscription ethereum.Subscription, eventChan <-chan *abi.AgentCollectionsManagerCollectionCreated) {
	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()

		listener.logger.Info("listening for AgentCollectionsManager.CollectionCreated events")

		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					return
				}
				if event == nil {
					continue
				}

				err := listener.collectionCreatedEventHandler.Handle(event)
				if err != nil {
					listener.logger.Error("failed to handle AgentCollectionsManager.CollectionCreated event", zap.Any("event", event), zap.Error(err))
				}

				for _, collectionListener := range listener.collectionListeners {
					err = collectionListener.Subscribe(event.CollectionAddress)
					if err != nil {
						listener.logger.Error("failed to notify collection listener about AgentCollectionsManager.CollectionCreated event", zap.Any("event", event), zap.Error(err))
					}
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
				case listener.errCh <- fmt.Errorf("agent collections manager collection created subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
