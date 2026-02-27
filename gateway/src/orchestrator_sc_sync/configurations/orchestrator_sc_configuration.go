package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_sc_sync/use_cases"
	"go.uber.org/dig"
)

func AddOrchestratorSCSyncConfiguration(container *dig.Container) {
	// Services
	err := container.Provide(
		services.NewPreviousAgentSignerImporter,
		dig.As(new(services.AgentPreviousSignerImporterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentSynchronizer,
		dig.As(new(services.AgentSynchronizerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentCollectionSynchronizer,
		dig.As(new(services.AgentCollectionSynchronizerInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Use cases
	err = container.Provide(
		use_cases.NewUpdateAgentOwner,
		dig.As(new(use_cases.UpdateAgentOwnerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewUpdateAgentAttribute,
		dig.As(new(use_cases.UpdateAgentAttributeInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewFinishAgentCreation,
		dig.As(new(use_cases.FinishAgentCreationInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewFinishAgentImportation,
		dig.As(new(use_cases.FinishAgentImportationInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewCreateAgentCollection,
		dig.As(new(use_cases.CreateAgentCollectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewPairProposalToAgent,
		dig.As(new(use_cases.PairProposalToAgentInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewRemoveAgentFromProposal,
		dig.As(new(use_cases.RemoveAgentFromProposalInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewUpsertAgentImage,
		dig.As(new(use_cases.UpsertAgentImageInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSynchronizeAllAgents,
		dig.As(new(use_cases.SynchronizeAllAgentsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSynchronizeAllAgentCollections,
		dig.As(new(use_cases.SynchronizeAllAgentCollectionsInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Event handlers
	err = container.Provide(
		event_handlers.NewMintedEventHandler,
		dig.As(new(event_handlers.MintedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewTransferEventHandler,
		dig.As(new(event_handlers.TransferEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewAgentSignerUpdatedEventHandler,
		dig.As(new(event_handlers.AgentSignerUpdatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewAgentAttributeUpdatedEventHandler,
		dig.As(new(event_handlers.AgentAttributeUpdatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewCollectionCreatedEventHandler,
		dig.As(new(event_handlers.CollectionCreatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewAgentImageUpdatedEventHandler,
		dig.As(new(event_handlers.AgentImageUpdatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewMintProposalCreatedEventHandler,
		dig.As(new(event_handlers.MintProposalCreatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewMintProposalRemovedEventHandler,
		dig.As(new(event_handlers.MintProposalRemovedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
