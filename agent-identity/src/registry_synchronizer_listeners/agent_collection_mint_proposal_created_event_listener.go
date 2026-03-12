package registry_synchronizer_listeners

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type AgentCollectionMintProposalCreatedListener struct {
	*AbstractUpdatableSubscriptionListener

	logger *zap.Logger

	agentCollectionRepository                repositories.AgentCollectionRepositoryInterface
	mintProposalCreatedEventHandler          interfaces.MintProposalCreatedEventHandlerInterface
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface

	trackedAddresses []common.Address
}

func NewAgentCollectionMintProposalCreatedListener(
	logger *zap.Logger,
	agentCollectionRepository repositories.AgentCollectionRepositoryInterface,
	mintProposalCreatedEventHandler interfaces.MintProposalCreatedEventHandlerInterface,
	agentCollectionEventSubscriptionProvider interfaces.AgentCollectionEventSubscriptionProviderInterface,
) *AgentCollectionMintProposalCreatedListener {
	return &AgentCollectionMintProposalCreatedListener{
		AbstractUpdatableSubscriptionListener: NewAbstractUpdatableSubscriptionListener(),

		logger: logger,

		agentCollectionRepository:                agentCollectionRepository,
		mintProposalCreatedEventHandler:          mintProposalCreatedEventHandler,
		agentCollectionEventSubscriptionProvider: agentCollectionEventSubscriptionProvider,

		trackedAddresses: []common.Address{},
	}
}

func (listener *AgentCollectionMintProposalCreatedListener) Listen() error {
	return listener.AbstractUpdatableSubscriptionListener.Listen()
}

func (listener *AgentCollectionMintProposalCreatedListener) SubscribeAll() error {
	collections, err := listener.agentCollectionRepository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get agent collections: %w", err)
	}

	collectionAddresses := make([]common.Address, 0, len(collections))
	for _, collection := range collections {
		collectionAddresses = append(collectionAddresses, common.HexToAddress(collection.Address))
	}

	if len(collectionAddresses) == 0 {
		listener.logger.Info("no agent collections found, skipping subscription to AgentCollection.MintProposalCreated events")
		return nil
	}

	listener.mu.Lock()
	listener.trackedAddresses = append([]common.Address(nil), collectionAddresses...)
	addressesToSubscribe := append([]common.Address(nil), collectionAddresses...)
	listener.mu.Unlock()

	events := make(chan *abi.AgentCollectionV1MintProposalCreated, 64)
	subscription, err := listener.agentCollectionEventSubscriptionProvider.StartMintProposalCreatedSubscription(addressesToSubscribe, events)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.MintProposalCreated events: %w", err)
	}

	listener.mu.Lock()
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info("subscribed to AgentCollection.MintProposalCreated events", zap.Int("collectionCount", len(addressesToSubscribe)))

	return nil
}

func (listener *AgentCollectionMintProposalCreatedListener) Subscribe(collectionAddress common.Address) error {
	listener.mu.Lock()
	for _, trackedAddress := range listener.trackedAddresses {
		if trackedAddress == collectionAddress {
			listener.mu.Unlock()
			return nil
		}
	}

	listener.trackedAddresses = append(listener.trackedAddresses, collectionAddress)
	addressesToSubscribe := append([]common.Address(nil), listener.trackedAddresses...)
	listener.mu.Unlock()

	events := make(chan *abi.AgentCollectionV1MintProposalCreated, 64)
	subscription, err := listener.agentCollectionEventSubscriptionProvider.StartMintProposalCreatedSubscription(addressesToSubscribe, events)
	if err != nil {
		return fmt.Errorf("failed to subscribe to AgentCollection.MintProposalCreated events for new address: %w", err)
	}

	listener.mu.Lock()
	if listener.subscription != nil {
		listener.oldSubscriptions = append(listener.oldSubscriptions, listener.subscription)
	}
	listener.subscription = subscription
	listener.mu.Unlock()

	listener.startSubscriptionHandler(subscription, events)
	listener.logger.Info("started new AgentCollection.MintProposalCreated subscription for collection", zap.String("collectionAddress", collectionAddress.Hex()))

	return nil
}

func (listener *AgentCollectionMintProposalCreatedListener) startSubscriptionHandler(subscription ethereum.Subscription, eventChan <-chan *abi.AgentCollectionV1MintProposalCreated) {
	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()

		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					return
				}
				if event == nil {
					continue
				}

				if err := listener.mintProposalCreatedEventHandler.Handle(event); err != nil {
					listener.logger.Error("failed to handle AgentCollection.MintProposalCreated event", zap.Any("event", event), zap.Error(err))
				}

				listener.mu.Lock()
				if event.Raw.BlockNumber > listener.lastBlockSeen {
					listener.lastBlockSeen = event.Raw.BlockNumber
				}
				listener.mu.Unlock()
			case err := <-subscription.Err():
				if err == nil {
					return
				}
				select {
				case listener.errCh <- fmt.Errorf("agent collection mint proposal created subscription error: %w", err):
				default:
				}
				return
			}
		}
	}()
}
