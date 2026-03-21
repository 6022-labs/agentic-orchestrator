package use_cases

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type SaveAgentPromptTemplateInterface interface {
	Save(request *requests.SaveAgentPromptTemplateRequest) error
}

type SaveAgentPromptTemplate struct {
	agentRegistryRequester       interfaces.AgentRegistryRequesterInterface
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface
}

func NewSaveAgentPromptTemplate(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface,
) *SaveAgentPromptTemplate {
	return &SaveAgentPromptTemplate{
		agentRegistryRequester:       agentRegistryRequester,
		agentPromptTemplateRequester: agentPromptTemplateRequester,
	}
}

func (s *SaveAgentPromptTemplate) Save(request *requests.SaveAgentPromptTemplateRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	agent, err := s.agentRegistryRequester.GetAgent(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
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
