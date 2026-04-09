package prompt_parameters

import (
	compute_interfaces "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/services/interfaces"
)

type PromptParameterBuilderFactory struct {
	agentRegistryRequester compute_interfaces.AgentRegistryRequesterInterface
}

func NewPromptParameterBuilderFactory(
	agentRegistryRequester compute_interfaces.AgentRegistryRequesterInterface,
) *PromptParameterBuilderFactory {
	return &PromptParameterBuilderFactory{
		agentRegistryRequester: agentRegistryRequester,
	}
}

func (ppbf *PromptParameterBuilderFactory) Create() compute_interfaces.PromptParametersBuilderInterface {
	return NewPromptParametersBuilder(ppbf.agentRegistryRequester)
}
