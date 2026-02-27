package orchestrator_sc_sync_event_listeners

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/subscribers"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AgentCollectionsManagerCollectionCreatedListener struct {
	logger                                             *zap.Logger
	collectionListeners                                []CollectionListenerInterface
	collectionCreatedEventHandler                      event_handlers.CollectionCreatedEventHandlerInterface
	agentCollectionsManagerCollectionCreatedSubscriber subscribers.AgentCollectionsManagerCollectionCreatedSubscriberInterface
}

type newAgentCollectionsManagerCollectionCreatedListenerParams struct {
	dig.In

	Logger                                             *zap.Logger
	CollectionListeners                                []CollectionListenerInterface `group:"collection_listeners"`
	CollectionCreatedEventHandler                      event_handlers.CollectionCreatedEventHandlerInterface
	AgentCollectionsManagerCollectionCreatedSubscriber subscribers.AgentCollectionsManagerCollectionCreatedSubscriberInterface
}

func NewAgentCollectionsManagerCollectionCreatedListener(
	params newAgentCollectionsManagerCollectionCreatedListenerParams,
) *AgentCollectionsManagerCollectionCreatedListener {
	return &AgentCollectionsManagerCollectionCreatedListener{
		logger:                        params.Logger,
		collectionListeners:           params.CollectionListeners,
		collectionCreatedEventHandler: params.CollectionCreatedEventHandler,
		agentCollectionsManagerCollectionCreatedSubscriber: params.AgentCollectionsManagerCollectionCreatedSubscriber,
	}
}

func (listener *AgentCollectionsManagerCollectionCreatedListener) Listen() error {
	ctx := context.Background()

	listener.logger.Debug("Subscribing to AgentCollectionsManager.CollectionCreated events")

	events := make(chan *abi.AgentCollectionsManagerCollectionCreated)

	subscription, err := listener.agentCollectionsManagerCollectionCreatedSubscriber.SubscribeCollectionCreated(ctx, events)
	if err != nil {
		return err
	}

	listener.logger.Info("Listening for AgentCollectionsManager.CollectionCreated events")

	for {
		select {
		case event := <-events:
			listener.logger.Info("Received AgentCollectionsManager.CollectionCreated event", zap.Any("event", event))

			err = listener.collectionCreatedEventHandler.Handle(event)
			if err != nil {
				listener.logger.Error("Failed to handle AgentCollectionsManager.CollectionCreated event", zap.Any("event", event), zap.Error(err))
			}

			for _, collectionListener := range listener.collectionListeners {
				err = collectionListener.Subscribe(event.CollectionAddress)
				if err != nil {
					listener.logger.Error("Failed to notify collection listener about AgentCollectionsManager.CollectionCreated event", zap.Any("event", event), zap.Error(err))
				}
			}
		case err := <-subscription.Err():
			listener.logger.Error("Subscription error", zap.Error(err))

			return err
		}
	}
}
