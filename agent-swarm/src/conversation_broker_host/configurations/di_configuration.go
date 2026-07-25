package configurations

import (
	conversation_broker_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/configurations"
	conversation_broker_blockchain_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/configurations"
	conversation_broker_db_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_db/configurations"
	host_settings "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_host/settings"
	conversation_broker_http_agent_identity_registry_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/configurations"
	conversation_broker_http_channels_and_integrations_bridge_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_channels_and_integrations_bridge/configurations"
	conversation_broker_http_remote_agent_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_remote_agent/configurations"
	conversation_broker_mvc_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_mvc/configurations"
	conversation_broker_redis_configuration "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_redis/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	conversation_broker_configuration.AddConversationBrokerConfiguration(container)
	conversation_broker_blockchain_configuration.AddConversationBrokerBlockchainConfiguration(container)
	conversation_broker_db_configuration.AddConversationBrokerDbConfiguration(container)
	conversation_broker_mvc_configuration.AddConversationBrokerMvcConfiguration(container)
	conversation_broker_redis_configuration.AddConversationBrokerRedisConfiguration(container)
	conversation_broker_http_remote_agent_configuration.AddConversationBrokerHttpRemoteAgentConfiguration(container)
	conversation_broker_http_agent_identity_registry_configuration.AddConversationBrokerHttpAgentIdentityRegistryConfiguration(container)
	conversation_broker_http_channels_and_integrations_bridge_configuration.AddConversationBrokerHttpChannelsAndIntegrationsBridgeConfiguration(container)

	return container
}
