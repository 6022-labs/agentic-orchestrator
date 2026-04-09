package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/compute"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/compute/prompt_parameters"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases"
	"go.uber.org/dig"
)

func AddComputeConfiguration(container *dig.Container) {
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
		dig.As(new(services.LLMChatCompletionUseCase)),
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

	// Internal LLM chat completion requester (bridges compute strategies to the LLM use case)
	err = container.Provide(
		services.NewInternalLLMConnectionChatCompletionRequester,
		dig.As(new(interfaces.LLMConnectionChatCompletionRequesterInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Prompt services
	err = container.Provide(
		services.NewPromptFromPromptTemplateParser,
		dig.As(new(interfaces.PromptFromPromptTemplateParserInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		prompt_parameters.NewPromptParameterBuilderFactory,
		dig.As(new(interfaces.PromptParametersBuilderFactoryInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Compute strategies
	err = container.Provide(
		compute.NewHumanComputeStrategy,
		dig.As(new(compute.HumanComputeStrategyInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		compute.NewCloneComputeStrategy,
		dig.As(new(compute.CloneComputeStrategyInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		compute.NewExpertComputeStrategy,
		dig.As(new(compute.ExpertComputeStrategyInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		compute.NewFacilitatorComputeStrategy,
		dig.As(new(compute.FacilitatorComputeStrategyInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Agent compute use case
	err = container.Provide(
		use_cases.NewAgentCompute,
		dig.As(new(use_cases.AgentComputeInterface)),
	)
	if err != nil {
		panic(err)
	}
}
