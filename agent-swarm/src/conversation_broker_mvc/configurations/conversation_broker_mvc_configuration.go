package configurations

import (
	common_metrics "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/metrics"
	common_metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/metrics/interfaces"
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/mvc/middlewares"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_mvc"
	"go.uber.org/dig"
)

func AddConversationBrokerMvcConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(func() common_metrics_interfaces.ApiRequestMetricsInterface {
		return common_metrics.NewApiRequestMetrics("agentic_orchestrator_conversation_broker_api")
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(mvc_middlewares.NewApiRequestMetricsMiddleware); err != nil {
		panic(err)
	}

	err := container.Provide(
		conversation_broker_mvc.NewBridgesController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewParticipantsController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewConversationsController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewPromptTemplatesController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		conversation_broker_mvc.NewRemoteAgentsController,
		dig.Group("controllers"),
		dig.As(new(conversation_broker_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
