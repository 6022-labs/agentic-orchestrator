package configurations

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"go.uber.org/dig"
)

func AddSlackBridgeConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewSlackSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewProcessingSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewOrchestratorSlackBridgeSettings)
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
		services.NewOrchestratorPrivateKeyFetcher,
		dig.As(new(interfaces.OrchestratorPrivateKeyFetcherInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewSlackThreadConclude,
		dig.As(new(interfaces.SlackThreadConcludeInterface)),
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
