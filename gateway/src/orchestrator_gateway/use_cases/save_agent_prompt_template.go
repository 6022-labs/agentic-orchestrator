package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway/use_cases/requests"
	"gorm.io/gorm"
)

type SaveAgentPromptTemplateInterface interface {
	Save(request *requests.SaveAgentPromptTemplateRequest) error
}

type SaveAgentPromptTemplate struct {
	agentRepository              repositories.AgentRepositoryInterface
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface
}

func NewSaveAgentPromptTemplate(
	agentRepository repositories.AgentRepositoryInterface,
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{
		agentRepository:              agentRepository,
		agentPromptTemplateRequester: agentPromptTemplateRequester,
	}
}

func (s *SaveAgentPromptTemplate) Save(request *requests.SaveAgentPromptTemplateRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agent, err := s.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
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
	if agent.Role == enums.AgentRoleHuman {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot save prompt template")
	}

	err = s.agentPromptTemplateRequester.SaveAgentPromptTemplate(agent.ID, request.Key, request.Content)
	if err != nil {
		return fmt.Errorf("failed to save agent prompt template: %w", err)
	}

	return nil
}
