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

type AgentCollectionMintProposalRemovedListener struct {
	logger                                       *zap.Logger
	mintProposalRemovedEventHandler              event_handlers.MintProposalRemovedEventHandlerInterface
	agentCollectionMintProposalRemovedSubscriber subscribers.AgentCollectionMintProposalRemovedSubscriberInterface

	errorChannel  chan error
	eventChannel  chan *abi.AgentCollectionV1MintProposalRemoved
	subscriptions []ethereum.Subscription
}

func NewAgentCollectionMintProposalRemovedListener(
	logger *zap.Logger,
	mintProposalRemovedEventHandler event_handlers.MintProposalRemovedEventHandlerInterface,
	agentCollectionMintProposalRemovedSubscriber subscribers.AgentCollectionMintProposalRemovedSubscriberInterface,
) *AgentCollectionMintProposalRemovedListener {
	return &AgentCollectionMintProposalRemovedListener{
		logger:                          logger,
		mintProposalRemovedEventHandler: mintProposalRemovedEventHandler,
		agentCollectionMintProposalRemovedSubscriber: agentCollectionMintProposalRemovedSubscriber,

		subscriptions: []ethereum.Subscription{},
		errorChannel:  make(chan error),
		eventChannel:  make(chan *abi.AgentCollectionV1MintProposalRemoved),
	}
}

func (listener *AgentCollectionMintProposalRemovedListener) Listen() error {
	for {
		select {
		case event := <-listener.eventChannel:
			listener.logger.Info("Received AgentCollection.MintProposalRemoved event", zap.Any("event", event))

			err := listener.mintProposalRemovedEventHandler.Handle(event)
			if err != nil {
				listener.logger.Error("Failed to handle AgentCollection.MintProposalRemoved event", zap.Any("event", event), zap.Error(err))
			}
		case err := <-listener.errorChannel:
			listener.logger.Error("Subscription error", zap.Error(err))

			return err
		}
	}
}

func (listener *AgentCollectionMintProposalRemovedListener) Subscribe(collectionAddress common.Address) error {
	ctx := context.Background()

	listener.logger.Debug("Subscribing to AgentCollection.MintProposalRemoved events", zap.String("collectionAddress", collectionAddress.Hex()))

	subscription, err := listener.agentCollectionMintProposalRemovedSubscriber.SubscribeMintProposalRemoved(ctx, collectionAddress, listener.eventChannel)
	if err != nil {
		return err
	}

	go func() {
		if err := <-subscription.Err(); err != nil {
			listener.errorChannel <- err
		}
	}()

	listener.logger.Info("Listening for AgentCollection.MintProposalRemoved events", zap.String("collectionAddress", collectionAddress.Hex()))
	listener.subscriptions = append(listener.subscriptions, subscription)

	return nil
}
