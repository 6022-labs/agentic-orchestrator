package configurations

import (
	gateway_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/configurations"
	gateway_host_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/settings"
	gateway_http_agent_identity_registry_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/configurations"
	gateway_http_agent_swarm_conversation_broker_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/configurations"
	gateway_http_channels_and_integrations_bridge_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/configurations"
	gateway_mvc_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc/configurations"
	gateway_redis_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(gateway_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	gateway_configurations.AddGatewayConfiguration(container)
	gateway_mvc_configurations.AddGatewayMvcConfiguration(container)
	gateway_redis_configurations.AddGatewayRedisConfiguration(container)
	gateway_http_channels_and_integrations_bridge_configurations.AddGatewayHttpChannelsAndIntegrationsBridgeConfiguration(container)
	gateway_http_agent_identity_registry_configurations.AddGatewayHttpAgentIdentityRegistryConfiguration(container)
	gateway_http_agent_swarm_conversation_broker_configurations.AddGatewayHttpAgentSwarmConversationBrokerConfiguration(container)

	return container
}
