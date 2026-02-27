package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/use_cases/requests"
)

type RemoveLLMConnectionInterface interface {
	Remove(request *requests.AgentIdRequest) error
}

type RemoveLLMConnection struct {
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface
}

func NewRemoveLLMConnection(
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface,
) *RemoveLLMConnection {
	return &RemoveLLMConnection{
		llmConnectionRepository: llmConnectionRepository,
	}
}

func (uc *RemoveLLMConnection) Remove(request *requests.AgentIdRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	err := uc.llmConnectionRepository.DeleteByAgentId(request.AgentId)
	if err != nil {
		return err
	}

	return nil
}
