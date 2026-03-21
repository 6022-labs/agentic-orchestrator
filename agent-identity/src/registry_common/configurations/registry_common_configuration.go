package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"go.uber.org/dig"
)

func AddRegistryCommonConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewBlockchainSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewAgentCollectionSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewOrchestratorSettings)
	if err != nil {
		panic(err)
	}

	err = container.Provide(settings.NewDecentralizedDomainNameProviderSettings)
	if err != nil {
		panic(err)
	}
}
