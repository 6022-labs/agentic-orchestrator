package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync_event_listeners"
	"go.uber.org/dig"
)

func AddOrchestratorScSyncEventListenersConfiguration(container *dig.Container) {
	// Listener
	err := container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionMintedListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionMintedListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionMintedListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionTransferListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionTransferListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionTransferListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionAgentSignerUpdatedListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionAgentSignerUpdatedListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionAgentSignerUpdatedListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionAgentAttributeUpdatedListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionAgentAttributeUpdatedListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionAgentAttributeUpdatedListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionAgentImageUpdatedListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionAgentImageUpdatedListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionAgentImageUpdatedListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionMintProposalCreatedListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionMintProposalCreatedListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionMintProposalCreatedListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(orchestrator_sc_sync_event_listeners.NewAgentCollectionMintProposalRemovedListener)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionMintProposalRemovedListener) orchestrator_sc_sync_event_listeners.ListenerInterface {
		return l
	},
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(l *orchestrator_sc_sync_event_listeners.AgentCollectionMintProposalRemovedListener) orchestrator_sc_sync_event_listeners.CollectionListenerInterface {
		return l
	},
		dig.Group("collection_listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.CollectionListenerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_sc_sync_event_listeners.NewAgentCollectionsManagerCollectionCreatedListener,
		dig.Group("listeners"),
		dig.As(new(orchestrator_sc_sync_event_listeners.ListenerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
