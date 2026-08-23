package registry_synchronizer_listeners

import (
	"context"
	"time"

	"fmt"
	"sync"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_listeners/metrics/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionsManagerCollectionCreatedListener struct {
	logger *zap.Logger

	chainEventMetrics metrics_interfaces.ChainEventMetricsInterface
	chainId           uint64

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
	chainEventMetrics metrics_interfaces.ChainEventMetricsInterface,
	blockchainSettings *common_settings.BlockchainSettings,
	agentCollectionsManagerSettings *settings.AgentCollectionsManagerSettings,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	agentCollectionsManagerEventSubscriptionProvider interfaces.AgentCollectionsManagerEventSubscriptionProviderInterface,
) *AgentCollectionsManagerCollectionCreatedListener {
	return &AgentCollectionsManagerCollectionCreatedListener{
		logger: logger,

		chainEventMetrics: chainEventMetrics,
		chainId:           uint64(blockchainSettings.ChainId),

		agentCollectionsManagerSettings:                  agentCollectionsManagerSettings,
		blockEventsQueue:                                 blockEventsQueue,
		agentCollectionsManagerEventSubscriptionProvider: agentCollectionsManagerEventSubscriptionProvider,

		errCh: make(chan error, 8),
	}
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) Listen(ctx context.Context) error {
	var err error
	select {
	case err = <-listener.errCh:
	case <-ctx.Done():
		// Graceful stop, but a subscription error already queued still wins.
		select {
		case err = <-listener.errCh:
		default:
		}
	}

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.subscription.Unsubscribe()
	}
	listener.mu.Unlock()

	listener.wg.Wait()

	return err
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) SubscribeAll(ctx context.Context) error {
	provider := listener.agentCollectionsManagerEventSubscriptionProvider
	events, subscription, err := provider.StartCollectionCreatedSubscription(ctx,
		[]common.Address{listener.agentCollectionsManagerSettings.SmartContractAddress},
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollectionsManager.CollectionCreated events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(ctx, subscription, events)

	listener.logger.Info(
		"subscribed to AgentCollectionsManager.CollectionCreated events",
		zap.String(
			"agentCollectionsManagerAddress",
			listener.agentCollectionsManagerSettings.SmartContractAddress.Hex(),
		),
	)

	return nil
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) startSubscriptionHandler(ctx context.Context,
	subscription ethereum.Subscription,
	eventChan <-chan *abi.AgentCollectionsManagerCollectionCreated,
) {
	listener.chainEventMetrics.RecordSubscriptionOpened(ctx, "CollectionCreated", listener.chainId)
	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()
		defer listener.chainEventMetrics.RecordSubscriptionClosed(ctx, "CollectionCreated", listener.chainId)

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

				start := time.Now()
				err := listener.blockEventsQueue.AddAgentCollectionsManagerCollectionCreatedEvents(
					ctx,
					event.Raw.BlockNumber,
					[]*abi.AgentCollectionsManagerCollectionCreated{event},
				)
				outcome := metrics_interfaces.ChainEventOutcomeHandled
				if err != nil {
					outcome = metrics_interfaces.ChainEventOutcomeFailed
					listener.logger.Error(
						"failed to enqueue AgentCollectionsManager.CollectionCreated event",
						zap.Uint64("blockNumber", event.Raw.BlockNumber),
						zap.Error(err),
					)
				}
				listener.chainEventMetrics.RecordEvent(
					ctx,
					"CollectionCreated",
					listener.chainId,
					outcome,
					time.Since(start),
				)
			case <-ctx.Done():
				return
			case err := <-subscription.Err():
				if err == nil {
					return
				}
				listener.chainEventMetrics.RecordSubscriptionError(ctx, "CollectionCreated", listener.chainId)
				select {
				case listener.errCh <- fmt.Errorf("agent collections manager collection created subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
