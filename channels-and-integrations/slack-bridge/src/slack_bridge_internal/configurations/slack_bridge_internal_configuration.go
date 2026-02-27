package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"go.uber.org/dig"
)

func AddSlackBridgeInternalConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(
		settings.NewProcessingSettings,
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewSlackThreadProcessor,
		dig.As(new(interfaces.SlackThreadProcessorInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackThreadProcessorDaemon,
		dig.As(new(interfaces.SlackThreadProcessorDaemonInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Use cases
	err = container.Provide(
		use_cases.NewSetWalletOwnershipToSlackUser,
		dig.As(new(use_cases.SetWalletOwnershipToSlackUserInterface)),
	)
	if err != nil {
		panic(err)
	}
}
