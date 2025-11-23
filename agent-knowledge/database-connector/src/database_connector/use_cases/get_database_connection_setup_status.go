package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/use_cases/responses"
)

type GetDatabaseConnectionSetupStatusInterface interface {
	Get(request *requests.AgentIdRequest) (*responses.GetDatabaseConnectionSetupStatusResponse, error)
}

type GetDatabaseConnectionSetupStatus struct {
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface
}

func NewGetDatabaseConnectionSetupStatus(
	databaseConnectionRepository repositories.DatabaseConnectionRepositoryInterface,
) *GetDatabaseConnectionSetupStatus {
	return &GetDatabaseConnectionSetupStatus{
		databaseConnectionRepository: databaseConnectionRepository,
	}
}

func (uc *GetDatabaseConnectionSetupStatus) Get(request *requests.AgentIdRequest) (*responses.GetDatabaseConnectionSetupStatusResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	status, err := uc.databaseConnectionRepository.GetSetupStatusByAgentId(request.AgentId)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return &responses.GetDatabaseConnectionSetupStatusResponse{
			IsCompleted: false,
			SetupStatus: nil,
		}, nil
	}

	return &responses.GetDatabaseConnectionSetupStatusResponse{
		IsCompleted: *status == enums.DatabaseConnectionSetupStatusEstablished,
		SetupStatus: status,
	}, nil
}
