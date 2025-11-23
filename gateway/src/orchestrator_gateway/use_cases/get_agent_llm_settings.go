package use_cases

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/repositories"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/use_cases/responses"
)

type GetAgentLLMSettingsInterface interface {
	Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.GetAgentLLMSettingsResponse, error)
}

type GetAgentLLMSettings struct {
	agentRepository        repositories.AgentRepositoryInterface
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface
}

func NewGetAgentLLMSettings(
	agentRepository repositories.AgentRepositoryInterface,
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface,
) *GetAgentLLMSettings {
	return &GetAgentLLMSettings{
		agentRepository:        agentRepository,
		llmConnectionRequester: llmConnectionRequester,
	}
}

func (g *GetAgentLLMSettings) Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.GetAgentLLMSettingsResponse, error) {
	err := request.EnsureAuthenticatedWallet()
	if err != nil {
		return nil, err
	}

	agent, err := g.agentRepository.GetByCollectionAddressAndCollectionTokenId(
		request.AgentCollectionAddress,
		request.AgentCollectionTokenId,
	)
	if err != nil {
		return nil, err
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
	if agent.Role == enums.AgentRoleHuman {
		return nil, errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot have LLM settings")
	}

	llmSettings, err := g.llmConnectionRequester.GetSettings(agent.ID)
	if err != nil {
		return nil, err
	}

	return &responses.GetAgentLLMSettingsResponse{
		Model:    llmSettings.Model,
		BaseUrl:  llmSettings.BaseUrl,
		Platform: llmSettings.Platform,
	}, nil
}
