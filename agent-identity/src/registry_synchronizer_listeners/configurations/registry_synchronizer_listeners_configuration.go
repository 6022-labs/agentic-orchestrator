package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_listeners"
	"go.uber.org/dig"
)

func AddRegistrySynchronizerListenersConfiguration(container *dig.Container) {
	err := addAgentCollectionMintedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionTransferListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionAgentSignerUpdatedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionAgentAttributeUpdatedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionAgentImageUpdatedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionMintProposalCreatedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionMintProposalRemovedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}

	err = addAgentCollectionsManagerCollectionCreatedListenerConfiguration(container)
	if err != nil {
		panic(err)
	}
}

func addAgentCollectionMintedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionMintedListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintedListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintedListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintedListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionTransferListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionTransferListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionTransferListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionTransferListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionTransferListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionAgentSignerUpdatedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionAgentSignerUpdatedListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentSignerUpdatedListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentSignerUpdatedListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentSignerUpdatedListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionAgentAttributeUpdatedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionAgentAttributeUpdatedListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentAttributeUpdatedListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentAttributeUpdatedListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentAttributeUpdatedListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionAgentImageUpdatedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionAgentImageUpdatedListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentImageUpdatedListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentImageUpdatedListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionAgentImageUpdatedListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionMintProposalCreatedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionMintProposalCreatedListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintProposalCreatedListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintProposalCreatedListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintProposalCreatedListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionMintProposalRemovedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(registry_synchronizer_listeners.NewAgentCollectionMintProposalRemovedListener)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintProposalRemovedListener) interfaces.CollectionEventSubscriberInterface {
		return l
	},
		dig.Group("collection_event_subscribers"),
		dig.As(new(interfaces.CollectionEventSubscriberInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintProposalRemovedListener) interfaces.OnBlockProcessedHandlerInterface {
		return l
	},
		dig.Group("on_block_processed_handlers"),
		dig.As(new(interfaces.OnBlockProcessedHandlerInterface)),
	)
	if err != nil {
		return err
	}

	err = container.Provide(func(l *registry_synchronizer_listeners.AgentCollectionMintProposalRemovedListener) registry_synchronizer_listeners.EventListenerInterface {
		return l
	},
		dig.Group("event_listeners"),
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
	)
	if err != nil {
		return err
	}

	return nil
}

func addAgentCollectionsManagerCollectionCreatedListenerConfiguration(container *dig.Container) error {
	err := container.Provide(
		registry_synchronizer_listeners.NewAgentCollectionsManagerCollectionCreatedListener,
		dig.As(new(registry_synchronizer_listeners.EventListenerInterface)),
		dig.Group("event_listeners"),
	)
	if err != nil {
		return err
	}

	return nil
}
