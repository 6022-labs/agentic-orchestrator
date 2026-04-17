package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"go.uber.org/dig"
)

func AddSlackBridgeConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewSlackSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewECDSASignature,
		dig.As(new(interfaces.ECDSASignatureInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewWalletOwnershipChallenge,
		dig.As(new(interfaces.WalletOwnershipChallengeInterface)),
	)
	if err != nil {
		panic(err)
	}
}
