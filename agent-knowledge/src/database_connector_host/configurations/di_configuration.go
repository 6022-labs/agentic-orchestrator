package configurations

import (
	database_connector_configuration "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/configurations"
	database_connector_db_agent_configuration "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_agent/configurations"
	database_connector_db_orchestrator_configuration "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_db_orchestrator/configurations"
	database_connector_host_settings "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_host/settings"
	database_connector_mvc_configuration "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_mvc/configurations"
	database_connector_vault_configuration "github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_vault/configurations"
	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(database_connector_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	database_connector_configuration.AddDatabaseConnectorConfiguration(container)
	database_connector_mvc_configuration.AddOrchestratorAgentDatabaseConnectorMvcConfiguration(container)
	database_connector_vault_configuration.AddDatabaseConnectorVaultConfiguration(container)
	database_connector_db_agent_configuration.AddDatabaseConnectorDbAgentConfiguration(container)
	database_connector_db_orchestrator_configuration.AddDatabaseConnectorDbOrchestratorConfiguration(container)

	return container
}
