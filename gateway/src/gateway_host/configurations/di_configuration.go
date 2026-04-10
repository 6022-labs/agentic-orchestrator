package configurations

import (
	gateway_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/configurations"
	gateway_blockchain_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_blockchain/configurations"
	gateway_host_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/settings"
	gateway_http_agent_identity_registry_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/configurations"
	gateway_http_agent_knowledge_database_connector_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_knowledge_database_connector/configurations"
	gateway_http_agent_reasoning_agent_compute_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_reasoning_agent_compute/configurations"
	gateway_http_platform_bridge_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_platform_bridge/configurations"
	gateway_mvc_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc/configurations"
	gateway_redis_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/configurations"
	gateway_vault_configurations "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_vault/configurations"

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
	gateway_vault_configurations.AddGatewayVaultConfiguration(container)
	gateway_redis_configurations.AddGatewayRedisConfiguration(container)
	gateway_blockchain_configurations.AddGatewayBlockchainConfiguration(container)
	gateway_http_platform_bridge_configurations.AddGatewayHttpPlatformBridgeConfiguration(container)
	gateway_http_agent_identity_registry_configurations.AddGatewayHttpAgentIdentityRegistryConfiguration(container)
	gateway_http_agent_reasoning_agent_compute_configurations.AddGatewayHttpAgentReasoningAgentComputeConfiguration(container)
	gateway_http_agent_knowledge_database_connector_configurations.AddGatewayHttpAgentKnowledgeDatabaseConnectorConfiguration(container)

	return container
}
