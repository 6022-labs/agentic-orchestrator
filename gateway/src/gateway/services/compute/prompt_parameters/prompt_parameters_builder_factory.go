package prompt_parameters

import (
	gateway_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
)

type PromptParameterBuilderFactory struct {
	agentRegistryRequester gateway_interfaces.AgentRegistryRequesterInterface
}

func NewPromptParameterBuilderFactory(
	agentRegistryRequester gateway_interfaces.AgentRegistryRequesterInterface,
) *PromptParameterBuilderFactory {
	return &PromptParameterBuilderFactory{
		agentRegistryRequester: agentRegistryRequester,
	}
}

func (ppbf *PromptParameterBuilderFactory) Create() gateway_interfaces.PromptParametersBuilderInterface {
	return NewPromptParametersBuilder(ppbf.agentRegistryRequester)
}
