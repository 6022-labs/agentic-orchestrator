package configurations

import (
	orchestrator_slack_bridge_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/configurations"
	slack_bridge_db_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db/configurations"
	orchestrator_slack_bridge_external_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/configurations"
	orchestrator_slack_bridge_external_host_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_host/settings"
	orchestrator_slack_bridge_external_mvc_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc/configurations"
	orchestrator_slack_bridge_http_orchestrator_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_orchestrator/configurations"
	orchestrator_slack_bridge_http_slack_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/configurations"
	orchestrator_slack_bridge_redis_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_redis/configurations"
	orchestrator_slack_bridge_vault_configurations "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_vault/configurations"

	"github.com/spf13/viper"
	"go.uber.org/dig"
)

func ConfigureDI() *dig.Container {
	container := dig.New()

	container.Provide(orchestrator_slack_bridge_external_host_settings.NewHostSettings)
	container.Provide(viper.GetViper)

	orchestrator_slack_bridge_configurations.AddSlackBridgeConfiguration(container)
	slack_bridge_db_configurations.AddSlackBridgeDbConfiguration(container)
	orchestrator_slack_bridge_vault_configurations.AddSlackBridgeVaultConfiguration(container)
	orchestrator_slack_bridge_redis_configurations.AddSlackBridgeRedisConfiguration(container)
	orchestrator_slack_bridge_external_configurations.AddSlackBridgeExternalConfiguration(container)
	orchestrator_slack_bridge_http_slack_configurations.AddSlackBridgeHttpSlackConfiguration(container)
	orchestrator_slack_bridge_external_mvc_configurations.AddSlackBridgeExternalMvcConfiguration(container)
	orchestrator_slack_bridge_http_orchestrator_configurations.AddSlackBridgeHttpOrchestratorConfiguration(container)

	return container
}
