package configurations

import (
	orchestrator_slack_bridge_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/configurations"
	slack_bridge_db_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_db/configurations"
	orchestrator_slack_bridge_http_slack_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/configurations"
	orchestrator_slack_bridge_internal_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal/configurations"
	orchestrator_slack_bridge_internal_host_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/settings"
	slack_bridge_internal_http_agent_reasoning_compute_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_http_agent_reasoning_compute/configurations"
	orchestrator_slack_bridge_internal_mvc_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_mvc/configurations"
	orchestrator_slack_bridge_redis_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/configurations"
	orchestrator_slack_bridge_vault_configurations "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_vault/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(orchestrator_slack_bridge_internal_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	orchestrator_slack_bridge_configurations.AddSlackBridgeConfiguration(container)
	slack_bridge_db_configurations.AddSlackBridgeDbConfiguration(container)
	orchestrator_slack_bridge_redis_configurations.AddSlackBridgeRedisConfiguration(container)
	orchestrator_slack_bridge_vault_configurations.AddSlackBridgeVaultConfiguration(container)
	orchestrator_slack_bridge_internal_configurations.AddSlackBridgeInternalConfiguration(container)
	orchestrator_slack_bridge_http_slack_configurations.AddSlackBridgeHttpSlackConfiguration(container)
	orchestrator_slack_bridge_internal_mvc_configurations.AddSlackBridgeInternalMvcConfiguration(container)
	slack_bridge_internal_http_agent_reasoning_compute_configurations.AddSlackBridgeInternalHttpAgentReasoningComputeConfiguration(container)

	return container
}
