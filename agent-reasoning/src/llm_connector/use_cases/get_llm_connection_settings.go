package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/use_cases/responses"
)

type GetLLMConnectionSettingsInterface interface {
	Get(request *requests.AgentIdRequest) (*responses.GetLLMConnectionSettingsResponse, error)
}

type GetLLMConnectionSettings struct {
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface
}

func NewGetLLMConnectionSettings(
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface,
) *GetLLMConnectionSettings {
	return &GetLLMConnectionSettings{
		llmConnectionRepository: llmConnectionRepository,
	}
}

func (uc *GetLLMConnectionSettings) Get(request *requests.AgentIdRequest) (*responses.GetLLMConnectionSettingsResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	llmConnection, err := uc.llmConnectionRepository.GetByAgentId(request.AgentId)
	if err != nil {
		return nil, err
	}
	if llmConnection == nil {
		return nil, fmt.Errorf("llm connection not found for agentId: %s", request.AgentId)
	}

	return &responses.GetLLMConnectionSettingsResponse{
		Model:    llmConnection.Model,
		BaseUrl:  llmConnection.BaseUrl,
		Platform: llmConnection.Platform,
	}, nil
}
