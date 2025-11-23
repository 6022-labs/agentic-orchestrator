package prompt_parameters

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"

type potentialAgentPromptParameter struct {
	*agentPromptParameter
	Name    string `json:"name"`
	IsAgent bool   `json:"isAgent"`
}

func newPotentialAgentPromptParameterFromModel(name string, agent *models.Agent) *potentialAgentPromptParameter {
	this := &potentialAgentPromptParameter{
		Name:    name,
		IsAgent: agent != nil,
	}

	if agent != nil {
		this.Name = agent.Name
		this.agentPromptParameter = newAgentPromptParameterFromModel(agent)
	}

	return this
}
