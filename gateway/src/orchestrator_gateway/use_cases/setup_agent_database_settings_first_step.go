package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
	"gorm.io/gorm"
)

type SetupAgentDatabaseSettingsFirstStepInterface interface {
	Setup(request *requests.SetupAgentDatabaseSettingsFirstStepRequest) (*responses.SetupAgentDatabaseSettingsFirstStepResponse, error)
}

type SetupAgentDatabaseSettingsFirstStep struct {
	agentRepository             repositories.AgentRepositoryInterface
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface
}

func NewSetupAgentDatabaseSettingsFirstStep(
	agentRepository repositories.AgentRepositoryInterface,
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface,
) *SetupAgentDatabaseSettingsFirstStep {
	return &SetupAgentDatabaseSettingsFirstStep{
		agentRepository:             agentRepository,
		databaseConnectionRequester: databaseConnectionRequester,
	}
}

func (u *SetupAgentDatabaseSettingsFirstStep) Setup(request *requests.SetupAgentDatabaseSettingsFirstStepRequest) (*responses.SetupAgentDatabaseSettingsFirstStepResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	agent, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	if agent == nil {
		return nil, errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if !agent.IsInternal {
		return nil, errors.NewValidationError("agentCollectionTokenId", "Agent is not internal")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return nil, errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}

	databaseConnectionResult, err := u.databaseConnectionRequester.Create(&interfaces.CreateDatabaseConnectionRequest{
		AgentId: agent.ID,
		Host:    request.Host,
		Port:    request.Port,
		Driver:  request.Driver,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return &responses.SetupAgentDatabaseSettingsFirstStepResponse{
		DatabaseName:     databaseConnectionResult.DatabaseName,
		DatabaseUsername: databaseConnectionResult.DatabaseUsername,
		DatabasePassword: databaseConnectionResult.DatabasePassword,
	}, nil
}
