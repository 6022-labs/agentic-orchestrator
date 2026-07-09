package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"go.uber.org/dig"
)

func AddRegistryCommonConfiguration(container *dig.Container) {
	// Settings
	err := container.Provide(settings.NewBlockchainSettings)
	if err != nil {
		panic(err)
	}



	err = container.Provide(settings.NewDecentralizedDomainNameProviderSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewDecentralizedDomainNameBuilder,
		dig.As(new(interfaces.DecentralizedDomainNameBuilderInterface)),
	)
	if err != nil {
		panic(err)
	}
}
