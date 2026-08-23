package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces"
	traces_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/traces/interfaces"
	"go.uber.org/dig"
)

func AddSlackBridgeConfiguration(container *dig.Container) {
	err := container.Provide(func() metrics_interfaces.ChannelMetricsInterface {
		return metrics.NewChannelMetrics(metrics_interfaces.PlatformSlack)
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() traces_interfaces.ChannelTracerInterface {
		return traces.NewChannelTracer(metrics_interfaces.PlatformSlack)
	})
	if err != nil {
		panic(err)
	}

	// Settings
	err = container.Provide(settings.NewSlackSettings)
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
