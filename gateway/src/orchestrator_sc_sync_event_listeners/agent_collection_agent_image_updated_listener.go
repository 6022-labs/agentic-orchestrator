package orchestrator_sc_sync_event_listeners

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/subscribers"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionAgentImageUpdatedListener struct {
	logger                                     *zap.Logger
	agentImageUpdatedEventHandler              event_handlers.AgentImageUpdatedEventHandlerInterface
	agentCollectionAgentImageUpdatedSubscriber subscribers.AgentCollectionAgentImageUpdatedSubscriberInterface

	errorChannel  chan error
	eventChannel  chan *abi.AgentCollectionV1AgentImageUpdated
	subscriptions []ethereum.Subscription
}

func NewAgentCollectionAgentImageUpdatedListener(
	logger *zap.Logger,
	agentImageUpdatedEventHandler event_handlers.AgentImageUpdatedEventHandlerInterface,
	agentCollectionAgentImageUpdatedSubscriber subscribers.AgentCollectionAgentImageUpdatedSubscriberInterface,
) *AgentCollectionAgentImageUpdatedListener {
	return &AgentCollectionAgentImageUpdatedListener{
		logger:                        logger,
		agentImageUpdatedEventHandler: agentImageUpdatedEventHandler,
		agentCollectionAgentImageUpdatedSubscriber: agentCollectionAgentImageUpdatedSubscriber,

		subscriptions: []ethereum.Subscription{},
		errorChannel:  make(chan error),
		eventChannel:  make(chan *abi.AgentCollectionV1AgentImageUpdated),
	}
}

func (listener *AgentCollectionAgentImageUpdatedListener) Listen() error {
	for {
		select {
		case event := <-listener.eventChannel:
			listener.logger.Info("Received AgentCollection.AgentImageUpdated event", zap.Any("event", event))

			err := listener.agentImageUpdatedEventHandler.Handle(event)
			if err != nil {
				listener.logger.Error("Failed to handle AgentCollection.AgentImageUpdated event", zap.Any("event", event), zap.Error(err))
			}
		case err := <-listener.errorChannel:
			listener.logger.Error("Subscription error", zap.Error(err))

			return err
		}
	}
}

func (listener *AgentCollectionAgentImageUpdatedListener) Subscribe(collectionAddress common.Address) error {
	ctx := context.Background()

	listener.logger.Debug("Subscribing to AgentCollection.AgentImageUpdated events", zap.String("collectionAddress", collectionAddress.Hex()))

	subscription, err := listener.agentCollectionAgentImageUpdatedSubscriber.SubscribeAgentImageUpdated(ctx, collectionAddress, listener.eventChannel)
	if err != nil {
		return err
	}

	go func() {
		if err := <-subscription.Err(); err != nil {
			listener.errorChannel <- err
		}
	}()

	listener.logger.Info("Listening for AgentCollection.AgentImageUpdated events", zap.String("collectionAddress", collectionAddress.Hex()))
	listener.subscriptions = append(listener.subscriptions, subscription)

	return nil
}
