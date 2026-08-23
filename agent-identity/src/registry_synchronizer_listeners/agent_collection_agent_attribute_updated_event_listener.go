package registry_synchronizer_listeners

import (
	"context"
	"time"

	"fmt"
	"slices"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	common_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_listeners/metrics/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionAgentAttributeUpdatedListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	chainEventMetrics metrics_interfaces.ChainEventMetricsInterface
	chainId           uint64

	agentCollectionRepository                repositories.AgentCollectionRepositoryInterface
	blockEventsQueue                         interfaces.BlockEventsQueueInterface
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface

	trackedAddresses []common.Address
}

func NewAgentCollectionAgentAttributeUpdatedListener(
	logger *zap.Logger,
	chainEventMetrics metrics_interfaces.ChainEventMetricsInterface,
	blockchainSettings *common_settings.BlockchainSettings,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface,
) *AgentCollectionAgentAttributeUpdatedListener {
	el := &AgentCollectionAgentAttributeUpdatedListener{
		logger: logger,

		chainEventMetrics: chainEventMetrics,
		chainId:           uint64(blockchainSettings.ChainId),

		agentCollectionRepository:                agentCollectionRepository,
		blockEventsQueue:                         blockEventsQueue,
		agentCollectionEventSubscriptionProvider: agentCollectionEventSubscriptionProvider,

		trackedAddresses: []common.Address{},
	}
	el.AbstractUpdatableSubscriptionListener = NewAbstractUpdatableSubscriptionListener(el.rebuildSubscription)

	return el
}

func (listener *AgentCollectionAgentAttributeUpdatedListener) Listen(ctx context.Context) error {
	return listener.AbstractUpdatableSubscriptionListener.Listen(ctx)
}

func (listener *AgentCollectionAgentAttributeUpdatedListener) SubscribeAll(ctx context.Context) error {
	collections, err := listener.agentCollectionRepository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get agent collections: %w", err)
	}

	collectionAddresses := make([]common.Address, 0, len(collections))
	for _, collection := range collections {
		collectionAddresses = append(collectionAddresses, common.HexToAddress(collection.Address))
	}

	if len(collectionAddresses) == 0 {
		listener.logger.Info(
			"no agent collections found, skipping subscription to AgentCollection.AgentAttributeUpdated events",
		)
		return nil
	}

	listener.mu.Lock()
	listener.trackedAddresses = append([]common.Address(nil), collectionAddresses...)
	addressesToSubscribe := append([]common.Address(nil), collectionAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentAttributeUpdatedSubscription(
		ctx,
		addressesToSubscribe,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.AgentAttributeUpdated events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(ctx, subscription, events)
	listener.logger.Info(
		"subscribed to AgentCollection.AgentAttributeUpdated events",
		zap.Int("collectionCount", len(addressesToSubscribe)),
	)

	return nil
}

func (listener *AgentCollectionAgentAttributeUpdatedListener) Subscribe(collectionAddress common.Address) error {
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

func (listener *AgentCollectionAgentAttributeUpdatedListener) rebuildSubscription(ctx context.Context) error {
	listener.mu.Lock()
	addressesToSubscribe := append([]common.Address(nil), listener.trackedAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentAttributeUpdatedSubscription(
		ctx,
		addressesToSubscribe,
	)
	if err != nil {
		return fmt.Errorf("failed to rebuild AgentCollection.AgentAttributeUpdated subscription: %w", err)
	}

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.oldSubscriptions = append(listener.oldSubscriptions, listener.subscription)
	}
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(ctx, subscription, events)
	listener.logger.Info(
		"rebuilt AgentCollection.AgentAttributeUpdated subscription",
		zap.Int("collectionCount", len(addressesToSubscribe)),
	)
	return nil
}

func (listener *AgentCollectionAgentAttributeUpdatedListener) startSubscriptionHandler(ctx context.Context,
	subscription ethereum.Subscription,
	eventChan <-chan *abi.AgentCollectionV1AgentAttributeUpdated,
) {
	listener.chainEventMetrics.RecordSubscriptionOpened(ctx, "AgentAttributeUpdated", listener.chainId)
	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()
		defer listener.chainEventMetrics.RecordSubscriptionClosed(ctx, "AgentAttributeUpdated", listener.chainId)

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
				err := listener.blockEventsQueue.AddAgentCollectionV1AgentAttributeUpdatedEvents(
					ctx,
					event.Raw.BlockNumber,
					[]*abi.AgentCollectionV1AgentAttributeUpdated{event},
				)
				outcome := metrics_interfaces.ChainEventOutcomeHandled
				if err != nil {
					outcome = metrics_interfaces.ChainEventOutcomeFailed
					listener.logger.Error(
						"failed to enqueue AgentCollection.AgentAttributeUpdated event",
						zap.Uint64("blockNumber", event.Raw.BlockNumber),
						zap.Error(err),
					)
				}
				listener.chainEventMetrics.RecordEvent(
					ctx,
					"AgentAttributeUpdated",
					listener.chainId,
					outcome,
					time.Since(start),
				)

				listener.onBlockSeen(event.Raw.BlockNumber)
			case <-ctx.Done():
				return
			case err := <-subscription.Err():
				if err == nil {
					return
				}
				listener.chainEventMetrics.RecordSubscriptionError(ctx, "AgentAttributeUpdated", listener.chainId)
				select {
				case listener.errCh <- fmt.Errorf("agent collection agent attribute updated subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
