package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/use_cases"
	"go.uber.org/dig"
)

func AddSlackBridgeInternalConfiguration(container *dig.Container) {
	// Use cases
	err := container.Provide(
		use_cases.NewSetWalletOwnershipToSlackUser,
		dig.As(new(use_cases.SetWalletOwnershipToSlackUserInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewPostBrokerCompletionToSlack,
		dig.As(new(use_cases.PostBrokerCompletionToSlackInterface)),
	)
	if err != nil {
		panic(err)
	}
}
