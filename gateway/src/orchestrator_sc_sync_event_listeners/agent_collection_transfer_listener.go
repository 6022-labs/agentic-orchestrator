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

type AgentCollectionTransferListener struct {
	logger                            *zap.Logger
	transferEventHandler              event_handlers.TransferEventHandlerInterface
	agentCollectionTransferSubscriber subscribers.AgentCollectionTransferSubscriberInterface

	errorChannel  chan error
	eventChannel  chan *abi.AgentCollectionV1Transfer
	subscriptions []ethereum.Subscription
}

func NewAgentCollectionTransferListener(
	logger *zap.Logger,
	transferEventHandler event_handlers.TransferEventHandlerInterface,
	agentCollectionTransferSubscriber subscribers.AgentCollectionTransferSubscriberInterface,
) *AgentCollectionTransferListener {
	return &AgentCollectionTransferListener{
		logger:                            logger,
		transferEventHandler:              transferEventHandler,
		agentCollectionTransferSubscriber: agentCollectionTransferSubscriber,

		subscriptions: []ethereum.Subscription{},
		errorChannel:  make(chan error),
		eventChannel:  make(chan *abi.AgentCollectionV1Transfer),
	}
}

func (listener *AgentCollectionTransferListener) Listen() error {
	for {
		select {
		case event := <-listener.eventChannel:
			listener.logger.Info("Received AgentCollection.Transfer event", zap.Any("event", event))

			err := listener.transferEventHandler.Handle(event)
			if err != nil {
				listener.logger.Error("Failed to handle AgentCollection.Transfer event", zap.Any("event", event), zap.Error(err))
			}
		case err := <-listener.errorChannel:
			listener.logger.Error("Subscription error", zap.Error(err))

			return err
		}
	}
}

func (listener *AgentCollectionTransferListener) Subscribe(collectionAddress common.Address) error {
	ctx := context.Background()

	listener.logger.Debug("Subscribing to AgentCollection.Transfer events", zap.String("collectionAddress", collectionAddress.Hex()))

	subscription, err := listener.agentCollectionTransferSubscriber.SubscribeTransfer(ctx, collectionAddress, listener.eventChannel)
	if err != nil {
		return err
	}

	go func() {
		if err := <-subscription.Err(); err != nil {
			listener.errorChannel <- err
		}
	}()

	listener.logger.Info("Listening for AgentCollection.Transfer events", zap.String("collectionAddress", collectionAddress.Hex()))
	listener.subscriptions = append(listener.subscriptions, subscription)

	return nil
}
