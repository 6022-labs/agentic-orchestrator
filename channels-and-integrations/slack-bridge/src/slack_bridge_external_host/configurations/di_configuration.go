package configurations

import (
	slack_bridge_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/configurations"
	slack_bridge_db_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/configurations"
	slack_bridge_external_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/configurations"
	slack_bridge_external_host_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/settings"
	slack_bridge_http_agent_identity_registry_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/configurations"
	slack_bridge_external_mvc_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc/configurations"
	slack_bridge_http_slack_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/configurations"
	slack_bridge_redis_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/configurations"
	slack_bridge_vault_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_vault/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(slack_bridge_external_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	slack_bridge_configurations.AddSlackBridgeConfiguration(container)
	slack_bridge_db_configurations.AddSlackBridgeDbConfiguration(container)
	slack_bridge_vault_configurations.AddSlackBridgeVaultConfiguration(container)
	slack_bridge_redis_configurations.AddSlackBridgeRedisConfiguration(container)
	slack_bridge_external_configurations.AddSlackBridgeExternalConfiguration(container)
	slack_bridge_http_slack_configurations.AddSlackBridgeHttpSlackConfiguration(container)
	slack_bridge_external_mvc_configurations.AddSlackBridgeExternalMvcConfiguration(container)
	slack_bridge_http_agent_identity_registry_configurations.AddSlackBridgeExternalHttpAgentIdentityRegistryConfiguration(container)

	return container
}
