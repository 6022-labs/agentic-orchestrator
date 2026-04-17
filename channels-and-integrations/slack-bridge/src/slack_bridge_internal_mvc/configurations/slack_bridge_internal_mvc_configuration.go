package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_mvc"
	"go.uber.org/dig"
)

func AddSlackBridgeInternalMvcConfiguration(container *dig.Container) {
	// Controllers
	err := container.Provide(
		slack_bridge_internal_mvc.NewWalletController,
		dig.Group("controllers"),
		dig.As(new(slack_bridge_internal_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		slack_bridge_internal_mvc.NewBrokerCallbackController,
		dig.Group("controllers"),
		dig.As(new(slack_bridge_internal_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
