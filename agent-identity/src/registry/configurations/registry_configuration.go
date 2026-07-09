package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/use_cases"
	"go.uber.org/dig"
)

func AddRegistryConfiguration(container *dig.Container) {
	err := container.Provide(
		use_cases.NewGetAgent,
		dig.As(new(use_cases.GetAgentInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListAgents,
		dig.As(new(use_cases.ListAgentsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListAgentCollections,
		dig.As(new(use_cases.ListAgentCollectionsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetMintedAgentByEnsDomain,
		dig.As(new(use_cases.GetMintedAgentByEnsDomainInterface)),
	)
	if err != nil {
		panic(err)
	}
}
