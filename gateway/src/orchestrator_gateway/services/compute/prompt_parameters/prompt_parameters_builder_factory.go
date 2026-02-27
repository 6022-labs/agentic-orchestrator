package prompt_parameters

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
)

type PromptParameterBuilderFactory struct {
	agentRepository repositories.AgentRepositoryInterface
}

func NewPromptParameterBuilderFactory(
	agentRepository repositories.AgentRepositoryInterface,
) *PromptParameterBuilderFactory {
	return &PromptParameterBuilderFactory{
		agentRepository: agentRepository,
	}
}

func (ppbf *PromptParameterBuilderFactory) Create() interfaces.PromptParametersBuilderInterface {
	return NewPromptParametersBuilder(ppbf.agentRepository)
}
