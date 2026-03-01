package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases"
	"go.uber.org/dig"
)

func AddLLMConnectorConfiguration(container *dig.Container) {
	// Use cases
	err := container.Provide(
		use_cases.NewSaveLLMConnection,
		dig.As(new(use_cases.SaveLLMConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewRemoveLLMConnection,
		dig.As(new(use_cases.RemoveLLMConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetLLMConnectionSettings,
		dig.As(new(use_cases.GetLLMConnectionSettingsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewLLMConnectionChatCompletion,
		dig.As(new(use_cases.LLMConnectionChatCompletionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetLLMConnectionSetupStatus,
		dig.As(new(use_cases.GetLLMConnectionSetupStatusInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewCreateMCPConnection,
		dig.As(new(use_cases.CreateMCPConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewRemoveMCPConnection,
		dig.As(new(use_cases.RemoveMCPConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewListMCPConnections,
		dig.As(new(use_cases.ListMCPConnectionsInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewAgentLLMSecretStorage,
		dig.As(new(services.AgentLLMSecretStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewMCPConnectionHeaderAuthValueStorage,
		dig.As(new(services.MCPConnectionHeaderAuthValueStorageInterface)),
	)
	if err != nil {
		panic(err)
	}
}
