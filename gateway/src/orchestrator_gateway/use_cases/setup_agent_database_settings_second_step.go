package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SetupAgentDatabaseSettingsSecondStepInterface interface {
	Setup(request *requests.AgentByCollectionTokenIdRequest) error
}

type SetupAgentDatabaseSettingsSecondStep struct {
	agentRepository             repositories.AgentRepositoryInterface
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface
}

func NewSetupAgentDatabaseSettingsSecondStep(
	agentRepository repositories.AgentRepositoryInterface,
	databaseConnectionRequester interfaces.DatabaseConnectionRequesterInterface,
) *SetupAgentDatabaseSettingsSecondStep {
	return &SetupAgentDatabaseSettingsSecondStep{
		agentRepository:             agentRepository,
		databaseConnectionRequester: databaseConnectionRequester,
	}
}

func (u *SetupAgentDatabaseSettingsSecondStep) Setup(request *requests.AgentByCollectionTokenIdRequest) error {
	err := request.EnsureAuthenticatedWallet()
	if err != nil {
		return err
	}

	agent, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get agent: %w", err)
	}
	if agent == nil {
		return errors.NewValidationError("agentCollectionTokenId", "agent not found")
	}
	if !agent.IsInternal {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is not internal")
	}
	if *agent.Owner != request.AuthenticatedWallet.Hex() {
		return errors.NewValidationError("agentCollectionTokenId", "You are not the owner of this agent")
	}

	copyFromAgentId := (*uuid.UUID)(nil)

	if agent.CloneOf != nil {
		agentToCopyFrom, err := u.agentRepository.GetByCollectionAddressAndCollectionTokenId(agent.CollectionAddress, *agent.CloneOf)
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to get agent to copy from: %w", err)
		}
		if agentToCopyFrom == nil {
			return fmt.Errorf("agent to copy from not found")
		}
		if !agentToCopyFrom.IsInternal {
			return fmt.Errorf("agent to copy from is not internal")
		}

		copyFromAgentId = &agentToCopyFrom.ID
	}

	err = u.databaseConnectionRequester.SetupAgentDatabase(agent.ID, copyFromAgentId)
	if err != nil {
		return fmt.Errorf("failed to setup agent database: %w", err)
	}

	return nil
}
