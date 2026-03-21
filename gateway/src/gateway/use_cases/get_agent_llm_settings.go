package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
)

type GetAgentLLMSettingsInterface interface {
	Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.GetAgentLLMSettingsResponse, error)
}

type GetAgentLLMSettings struct {
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface
}

func NewGetAgentLLMSettings(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	llmConnectionRequester interfaces.LLMConnectionRequesterInterface,
) *GetAgentLLMSettings {
	return &GetAgentLLMSettings{
		agentRegistryRequester: agentRegistryRequester,
		llmConnectionRequester: llmConnectionRequester,
	}
}

func (g *GetAgentLLMSettings) Get(request *requests.AgentByCollectionTokenIdRequest) (*responses.GetAgentLLMSettingsResponse, error) {
	err := request.EnsureAuthenticatedWallet()
	if err != nil {
		return nil, err
	}

	agent, err := g.agentRegistryRequester.GetAgent(
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
