package use_cases

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
)

type GetAgentPromptTemplateInterface interface {
	Get(request *requests.GetAgentPromptTemplateRequest) (*string, error)
}

type GetAgentPromptTemplate struct {
	agentRegistryRequester       interfaces.AgentRegistryRequesterInterface
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface
}

func NewGetAgentPromptTemplate(
	agentRegistryRequester interfaces.AgentRegistryRequesterInterface,
	agentPromptTemplateRequester interfaces.AgentPromptTemplateRequesterInterface,
) *GetAgentPromptTemplate {
	return &GetAgentPromptTemplate{
		agentRegistryRequester:       agentRegistryRequester,
		agentPromptTemplateRequester: agentPromptTemplateRequester,
	}
}

func (g *GetAgentPromptTemplate) Get(request *requests.GetAgentPromptTemplateRequest) (*string, error) {
	err := request.Validate()
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
		return nil, errors.NewValidationError("agentCollectionTokenId", "Agent is human, cannot have prompt templates")
	}

	promptTemplate, err := g.agentPromptTemplateRequester.GetAgentPromptTemplateByAgentIdAndKey(agent.ID, request.Key)
	if err != nil {
		return nil, err
	}

	return promptTemplate, nil
}
