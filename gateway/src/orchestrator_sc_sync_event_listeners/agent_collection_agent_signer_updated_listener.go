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

type AgentCollectionAgentSignerUpdatedListener struct {
	logger                                      *zap.Logger
	agentSignerUpdatedEventHandler              event_handlers.AgentSignerUpdatedEventHandlerInterface
	agentCollectionAgentSignerUpdatedSubscriber subscribers.AgentCollectionAgentSignerUpdatedSubscriberInterface

	errorChannel  chan error
	eventChannel  chan *abi.AgentCollectionV1AgentSignerUpdated
	subscriptions []ethereum.Subscription
}

func NewAgentCollectionAgentSignerUpdatedListener(
	logger *zap.Logger,
	agentSignerUpdatedEventHandler event_handlers.AgentSignerUpdatedEventHandlerInterface,
	agentCollectionAgentSignerUpdatedSubscriber subscribers.AgentCollectionAgentSignerUpdatedSubscriberInterface,
) *AgentCollectionAgentSignerUpdatedListener {
	return &AgentCollectionAgentSignerUpdatedListener{
		logger:                         logger,
		agentSignerUpdatedEventHandler: agentSignerUpdatedEventHandler,
		agentCollectionAgentSignerUpdatedSubscriber: agentCollectionAgentSignerUpdatedSubscriber,

		subscriptions: []ethereum.Subscription{},
		errorChannel:  make(chan error),
		eventChannel:  make(chan *abi.AgentCollectionV1AgentSignerUpdated),
	}
}

func (listener *AgentCollectionAgentSignerUpdatedListener) Listen() error {
	for {
		select {
		case event := <-listener.eventChannel:
			listener.logger.Info("Received AgentCollection.AgentSignerUpdated event", zap.Any("event", event))

			err := listener.agentSignerUpdatedEventHandler.Handle(event)
			if err != nil {
				listener.logger.Error("Failed to handle AgentCollection.AgentSignerUpdated event", zap.Any("event", event), zap.Error(err))
			}
		case err := <-listener.errorChannel:
			listener.logger.Error("Subscription error", zap.Error(err))

			return err
		}
	}
}

func (listener *AgentCollectionAgentSignerUpdatedListener) Subscribe(collectionAddress common.Address) error {
	ctx := context.Background()

	listener.logger.Debug("Subscribing to AgentCollection.AgentSignerUpdated events", zap.String("collectionAddress", collectionAddress.Hex()))

	subscription, err := listener.agentCollectionAgentSignerUpdatedSubscriber.SubscribeAgentSignerUpdated(ctx, collectionAddress, listener.eventChannel)
	if err != nil {
		return err
	}

	go func() {
		if err := <-subscription.Err(); err != nil {
			listener.errorChannel <- err
		}
	}()

	listener.logger.Info("Listening for AgentCollection.AgentSignerUpdated events", zap.String("collectionAddress", collectionAddress.Hex()))
	listener.subscriptions = append(listener.subscriptions, subscription)

	return nil
}
