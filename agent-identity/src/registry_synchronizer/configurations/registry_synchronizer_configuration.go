package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/services/interfaces"
	"go.uber.org/dig"
)

func AddRegistrySynchronizerConfiguration(container *dig.Container) {
	// Services
	err := container.Provide(
		services.NewUniswapV3BlockEventsQueue,
		dig.As(new(interfaces.BlockEventsQueueInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewBlockEventsProcessor,
		dig.As(new(interfaces.BlockEventsProcessorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewBlockEventsGapEnqueuer,
		dig.As(new(interfaces.BlockEventsGapEnqueuerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewBlockEventsDequeueDaemon,
		dig.As(new(interfaces.BlockEventsDequeueDaemonInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewCollectionCreatedEventHandler,
		dig.As(new(interfaces.CollectionCreatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewMintedEventHandler,
		dig.As(new(interfaces.MintedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewTransferEventHandler,
		dig.As(new(interfaces.TransferEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentSignerUpdatedEventHandler,
		dig.As(new(interfaces.AgentSignerUpdatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentAttributeUpdatedEventHandler,
		dig.As(new(interfaces.AgentAttributeUpdatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentImageUpdatedEventHandler,
		dig.As(new(interfaces.AgentImageUpdatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewMintProposalCreatedEventHandler,
		dig.As(new(interfaces.MintProposalCreatedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewMintProposalRemovedEventHandler,
		dig.As(new(interfaces.MintProposalRemovedEventHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
