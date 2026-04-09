package configurations

import (
	registry_common_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/configurations"
	registry_common_db_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/configurations"
	registry_common_vault_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_vault/configurations"
	registry_synchronizer_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer/configurations"
	registry_synchronizer_blockchain_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/configurations"
	registry_synchronizer_host_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_host/settings"
	registry_synchronizer_listeners_configurations "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_listeners/configurations"

	"github.com/knadh/koanf/v2"
	"go.uber.org/dig"
)

func ConfigureDI(config *koanf.Koanf) *dig.Container {
	container := dig.New()

	container.Provide(registry_synchronizer_host_settings.NewHostSettings)
	container.Provide(func() *koanf.Koanf {
		return config
	})

	registry_common_configurations.AddRegistryCommonConfiguration(container)
	registry_common_db_configurations.AddRegistryCommonDbConfiguration(container)
	registry_common_vault_configurations.AddRegistryCommonVaultConfiguration(container)
	registry_synchronizer_configurations.AddRegistrySynchronizerConfiguration(container)
	registry_synchronizer_listeners_configurations.AddRegistrySynchronizerListenersConfiguration(container)
	registry_synchronizer_blockchain_configurations.AddRegistrySynchronizerBlockchainConfiguration(container)

	return container
}
