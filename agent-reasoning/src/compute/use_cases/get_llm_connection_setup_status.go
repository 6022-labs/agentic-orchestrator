package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/use_cases/responses"
)

type GetLLMConnectionSetupStatusInterface interface {
	Get(request *requests.AgentIdRequest) (*responses.GetLLMConnectionSetupStatusResponse, error)
}

type GetLLMConnectionSetupStatus struct {
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface
}

func NewGetLLMConnectionSetupStatus(
	llmConnectionRepository repositories.LLMConnectionRepositoryInterface,
) *GetLLMConnectionSetupStatus {
	return &GetLLMConnectionSetupStatus{
		llmConnectionRepository: llmConnectionRepository,
	}
}

func (uc *GetLLMConnectionSetupStatus) Get(request *requests.AgentIdRequest) (*responses.GetLLMConnectionSetupStatusResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	status, err := uc.llmConnectionRepository.GetSetupStatusByAgentId(request.AgentId)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return &responses.GetLLMConnectionSetupStatusResponse{
			IsCompleted: false,
			SetupStatus: nil,
		}, nil
	}

	return &responses.GetLLMConnectionSetupStatusResponse{
		IsCompleted: *status == enums.LLMConnectionSetupStatusEstablished,
		SetupStatus: status,
	}, nil
}
