package orchestrator_sc_sync_event_listeners

import (
	"context"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/abi"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/event_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync/subscribers"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionAgentAttributeUpdatedListener struct {
	logger                                         *zap.Logger
	agentAttributeUpdatedEventHandler              event_handlers.AgentAttributeUpdatedEventHandlerInterface
	agentCollectionAgentAttributeUpdatedSubscriber subscribers.AgentCollectionAgentAttributeUpdatedSubscriberInterface

	errorChannel  chan error
	eventChannel  chan *abi.AgentCollectionV1AgentAttributeUpdated
	subscriptions []ethereum.Subscription
}

func NewAgentCollectionAgentAttributeUpdatedListener(
	logger *zap.Logger,
	agentAttributeUpdatedEventHandler event_handlers.AgentAttributeUpdatedEventHandlerInterface,
	agentCollectionAgentAttributeUpdatedSubscriber subscribers.AgentCollectionAgentAttributeUpdatedSubscriberInterface,
) *AgentCollectionAgentAttributeUpdatedListener {
	return &AgentCollectionAgentAttributeUpdatedListener{
		logger:                            logger,
		agentAttributeUpdatedEventHandler: agentAttributeUpdatedEventHandler,
		agentCollectionAgentAttributeUpdatedSubscriber: agentCollectionAgentAttributeUpdatedSubscriber,

		subscriptions: []ethereum.Subscription{},
		errorChannel:  make(chan error),
		eventChannel:  make(chan *abi.AgentCollectionV1AgentAttributeUpdated),
	}
}

func (listener *AgentCollectionAgentAttributeUpdatedListener) Listen() error {
	for {
		select {
		case event := <-listener.eventChannel:
			listener.logger.Info("Received AgentCollection.AgentAttributeUpdated event", zap.Any("event", event))

			err := listener.agentAttributeUpdatedEventHandler.Handle(event)
			if err != nil {
				listener.logger.Error("Failed to handle AgentCollection.AgentAttributeUpdated event", zap.Any("event", event), zap.Error(err))
			}
		case err := <-listener.errorChannel:
			listener.logger.Error("Subscription error", zap.Error(err))

			return err
		}
	}
}

func (listener *AgentCollectionAgentAttributeUpdatedListener) Subscribe(collectionAddress common.Address) error {
	ctx := context.Background()

	listener.logger.Debug("Subscribing to AgentCollection.AgentAttributeUpdated events", zap.String("collectionAddress", collectionAddress.Hex()))

	subscription, err := listener.agentCollectionAgentAttributeUpdatedSubscriber.SubscribeAgentAttributeUpdated(ctx, collectionAddress, listener.eventChannel)
	if err != nil {
		return err
	}

	go func() {
		if err := <-subscription.Err(); err != nil {
			listener.errorChannel <- err
		}
	}()

	listener.logger.Info("Listening for AgentCollection.AgentAttributeUpdated events", zap.String("collectionAddress", collectionAddress.Hex()))
	listener.subscriptions = append(listener.subscriptions, subscription)

	return nil
}
