package configurations

import (
	"net/http"
	"time"

	service_interfaces "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_remote_agent/settings"
	"go.uber.org/dig"
)

func AddConversationBrokerHttpRemoteAgentConfiguration(container *dig.Container) {
	err := container.Provide(settings.NewRemoteAgentSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(remoteAgentSettings *settings.RemoteAgentSettings) *http.Client {
		return &http.Client{Timeout: time.Duration(remoteAgentSettings.RequestTimeoutSeconds) * time.Second}
	}, dig.Name("agent-node-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(func(remoteAgentSettings *settings.RemoteAgentSettings) *http.Client {
		return &http.Client{Timeout: time.Duration(remoteAgentSettings.WellKnownRequestTimeoutSeconds) * time.Second}
	}, dig.Name("agent-node-well-known-http-client"))
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		clients.NewRemoteAgentClient,
		dig.As(new(clients.RemoteAgentClientInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRemoteAgentCardRequester,
		dig.As(new(service_interfaces.RemoteAgentCardRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewRemoteAgentCompletionsRequester,
		dig.As(new(service_interfaces.RemoteAgentCompletionsRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}
}
