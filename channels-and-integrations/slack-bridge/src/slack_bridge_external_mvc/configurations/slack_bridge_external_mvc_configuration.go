package configurations

import (
	orchestrator_slack_bridge_external_mvc "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc/middlewares"
	"go.uber.org/dig"
)

func AddSlackBridgeExternalMvcConfiguration(container *dig.Container) {
	// Controllers
	err := container.Provide(
		orchestrator_slack_bridge_external_mvc.NewEventHandlerController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_slack_bridge_external_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		orchestrator_slack_bridge_external_mvc.NewCommandHandlerController,
		dig.Group("controllers"),
		dig.As(new(orchestrator_slack_bridge_external_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Middlewares
	err = container.Provide(
		middlewares.NewSlackSignatureMiddleware,
	)
	if err != nil {
		panic(err)
	}
}
