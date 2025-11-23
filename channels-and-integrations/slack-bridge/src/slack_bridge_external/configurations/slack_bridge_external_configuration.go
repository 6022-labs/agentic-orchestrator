package configurations

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"go.uber.org/dig"
)

func AddSlackBridgeExternalConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewCacheSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewVerifyWalletOwnershipSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewSlackChannelSynchronizer,
		dig.As(new(interfaces.SlackChannelSynchronizerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackThreadMessagesAdder,
		dig.As(new(interfaces.SlackThreadMessagesAdderInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackChannelSynchronizerCached,
		dig.As(new(interfaces.SlackChannelSynchronizerCachedInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackAgentSynchronizer,
		dig.As(new(interfaces.SlackAgentSynchronizerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackUserSynchronizer,
		dig.As(new(interfaces.SlackUserSynchronizerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackUserSynchronizerCached,
		dig.As(new(interfaces.SlackUserSynchronizerCachedInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Event handlers
	err = container.Provide(
		event_handlers.NewMessageCreatedHandler,
		dig.As(new(event_handlers.MessageCreatedHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewMessageChangedHandler,
		dig.As(new(event_handlers.MessageChangedHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewMessageDeletedHandler,
		dig.As(new(event_handlers.MessageDeletedHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewChannelLeftHandler,
		dig.As(new(event_handlers.ChannelLeftHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewChannelArchiveHandler,
		dig.As(new(event_handlers.ChannelArchiveHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewChannelCreatedHandler,
		dig.As(new(event_handlers.ChannelCreatedHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewChannelDeletedHandler,
		dig.As(new(event_handlers.ChannelDeletedHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		event_handlers.NewChannelIdChangedHandler,
		dig.As(new(event_handlers.ChannelIdChangedHandlerInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Command handlers
	err = container.Provide(
		command_handlers.NewAddAgentToChannel,
		dig.As(new(command_handlers.AddAgentToChannelInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		command_handlers.NewRemoveAgentFromChannel,
		dig.As(new(command_handlers.RemoveAgentFromChannelInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		command_handlers.NewListAgentInChannel,
		dig.As(new(command_handlers.ListAgentInChannelInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		command_handlers.NewVerifyWalletOwnership,
		dig.As(new(command_handlers.VerifyWalletOwnershipInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		command_handlers.NewPairSlackUserToHumanAgent,
		dig.As(new(command_handlers.PairSlackUserToHumanAgentInterface)),
	)
	if err != nil {
		panic(err)
	}
}
