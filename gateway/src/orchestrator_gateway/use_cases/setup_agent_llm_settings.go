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

type SetupAgentLLMSettingsInterface interface {
	Setup(request *requests.SetupAgentLLMSettingsRequest) error
}

type SetupAgentLLMSettings struct {
	agentRepository        repositories.AgentRepositoryInterface
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface
}

func NewSetupAgentLLMSettings(
	agentRepository repositories.AgentRepositoryInterface,
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface,
) *SetupAgentLLMSettings {
	return &SetupAgentLLMSettings{
		agentRepository:        agentRepository,
		llmConnectionRequester: llmConnectionRequester,
	}
}

func (u *SetupAgentLLMSettings) Setup(request *requests.SetupAgentLLMSettingsRequest) error {
	err := request.Validate()
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
	if agent.Role == enums.AgentRoleHuman {
		return errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot setup LLM settings")
	}

	err = u.llmConnectionRequester.Save(&interfaces.SaveLLMConnectionRequest{
		AgentId:  agent.ID,
		Model:    request.Model,
		ApiKey:   request.ApiKey,
		BaseUrl:  request.BaseUrl,
		Platform: request.Platform,

		AmazonPlatformRegion: request.AmazonPlatformRegion,
	})
	if err != nil {
		return fmt.Errorf("failed to save LLM connection: %w", err)
	}

	return nil
}
