package registry_synchronizer_listeners

import (
	"fmt"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionsManagerCollectionCreatedListener struct {
	logger *zap.Logger

	agentCollectionsManagerSettings                  *settings.AgentCollectionsManagerSettings
	blockEventsQueue                                 interfaces.BlockEventsQueueInterface
	agentCollectionsManagerEventSubscriptionProvider interfaces.AgentCollectionsManagerEventSubscriptionProviderInterface

	subscription ethereum.Subscription
	mu           sync.Mutex
	errCh        chan error
	wg           sync.WaitGroup
}

func NewAgentCollectionsManagerCollectionCreatedListener(
	logger *zap.Logger,
	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	agentCollectionsManagerEventSubscriptionProvider interfaces.AgentCollectionsManagerEventSubscriptionProviderInterface,
) *AgentCollectionsManagerCollectionCreatedListener {
	return &AgentCollectionsManagerCollectionCreatedListener{
		logger: logger,

		agentCollectionsManagerSettings:                  agentCollectionsManagerSettings,
		blockEventsQueue:                                 blockEventsQueue,
		agentCollectionsManagerEventSubscriptionProvider: agentCollectionsManagerEventSubscriptionProvider,

		errCh: make(chan error, 8),
	}
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) Listen() error {
	err := <-listener.errCh

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.subscription.Unsubscribe()
	}
	listener.mu.Unlock()

	listener.wg.Wait()

	return err
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) SubscribeAll() error {
	events, subscription, err := listener.agentCollectionsManagerEventSubscriptionProvider.StartCollectionCreatedSubscription(
		[]common.Address{listener.agentCollectionsManagerSettings.SmartContractAddress},
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

				if err := listener.blockEventsQueue.AddAgentCollectionsManagerCollectionCreatedEvents(event.Raw.BlockNumber, []*abi.AgentCollectionsManagerCollectionCreated{event}); err != nil {
					listener.logger.Error(
						"failed to enqueue AgentCollectionsManager.CollectionCreated event",
						zap.Uint64("blockNumber", event.Raw.BlockNumber),
						zap.Error(err),
					)
				}
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
