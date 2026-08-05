package configurations

import (
	registry_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/configurations"
	registry_common_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/configurations"
	registry_common_db_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/configurations"
	registry_host_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/settings"
	registry_mvc_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_mvc/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(registry_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	registry_configurations.AddRegistryConfiguration(container)
	registry_mvc_configurations.AddRegistryMvcConfiguration(container)
	registry_common_configurations.AddRegistryCommonConfiguration(container)
	registry_common_db_configurations.AddRegistryCommonDbConfiguration(container)

	return container
}
