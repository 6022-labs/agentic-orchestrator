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

type AgentCollectionAgentAddressAddedListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	chainEventMetrics metrics_interfaces.ChainEventMetricsInterface
	chainId           uint64

	agentCollectionRepository                repositories.AgentCollectionRepositoryInterface
	blockEventsQueue                         interfaces.BlockEventsQueueInterface
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface

	trackedAddresses []common.Address
}

func NewAgentCollectionAgentAddressAddedListener(
	logger *zap.Logger,
	chainEventMetrics metrics_interfaces.ChainEventMetricsInterface,
	blockchainSettings *common_settings.BlockchainSettings,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	blockEventsQueue interfaces.BlockEventsQueueInterface,
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface,
) *AgentCollectionAgentAddressAddedListener {
	el := &AgentCollectionAgentAddressAddedListener{
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

func (listener *AgentCollectionAgentAddressAddedListener) Listen(ctx context.Context) error {
	return listener.AbstractUpdatableSubscriptionListener.Listen(ctx)
}

func (listener *AgentCollectionAgentAddressAddedListener) SubscribeAll(ctx context.Context) error {
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
			"no agent collections found, skipping subscription to AgentCollection.AgentAddressAdded events",
		)
		return nil
	}

	listener.mu.Lock()
	listener.trackedAddresses = append([]common.Address(nil), collectionAddresses...)
	addressesToSubscribe := append([]common.Address(nil), collectionAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentAddressAddedSubscription(
		ctx,
		addressesToSubscribe,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.AgentAddressAdded events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(ctx, subscription, events)
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

func (listener *AgentCollectionAgentAddressAddedListener) rebuildSubscription(ctx context.Context) error {
	listener.mu.Lock()
	addressesToSubscribe := append([]common.Address(nil), listener.trackedAddresses...)
	listener.mu.Unlock()

	events, subscription, err := listener.agentCollectionEventSubscriptionProvider.StartAgentAddressAddedSubscription(
		ctx,
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

	listener.startSubscriptionHandler(ctx, subscription, events)
	listener.logger.Info(
		"rebuilt AgentCollection.AgentAddressAdded subscription",
		zap.Int("collectionCount", len(addressesToSubscribe)),
	)
	return nil
}

func (listener *AgentCollectionAgentAddressAddedListener) startSubscriptionHandler(ctx context.Context,
	subscription ethereum.Subscription,
	eventChan <-chan *abi.AgentCollectionV1AgentAddressAdded,
) {
	listener.chainEventMetrics.RecordSubscriptionOpened(ctx, "AgentAddressAdded", listener.chainId)
	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()
		defer listener.chainEventMetrics.RecordSubscriptionClosed(ctx, "AgentAddressAdded", listener.chainId)

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
				err := listener.blockEventsQueue.AddAgentCollectionV1AgentAddressAddedEvents(
					ctx,
					event.Raw.BlockNumber,
					[]*abi.AgentCollectionV1AgentAddressAdded{event},
				)
				outcome := metrics_interfaces.ChainEventOutcomeHandled
				if err != nil {
					outcome = metrics_interfaces.ChainEventOutcomeFailed
					listener.logger.Error(
						"failed to enqueue AgentCollection.AgentAddressAdded event",
						zap.Uint64("blockNumber", event.Raw.BlockNumber),
						zap.Error(err),
					)
				}
				listener.chainEventMetrics.RecordEvent(
					ctx,
					"AgentAddressAdded",
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
				listener.chainEventMetrics.RecordSubscriptionError(ctx, "AgentAddressAdded", listener.chainId)
				select {
				case listener.errCh <- fmt.Errorf("agent collection agent address added subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
