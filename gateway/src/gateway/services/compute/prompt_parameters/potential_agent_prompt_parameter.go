package prompt_parameters

import "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"

type potentialAgentPromptParameter struct {
	*agentPromptParameter
	Name    string `json:"name"`
	IsAgent bool   `json:"isAgent"`
}

func newPotentialAgentPromptParameterFromResponse(name string, agent *responses.AgentResponse) *potentialAgentPromptParameter {
	this := &potentialAgentPromptParameter{
		Name:    name,
		IsAgent: agent != nil,
	}

	if agent != nil {
		this.Name = agent.Name
		this.agentPromptParameter = newAgentPromptParameterFromResponse(agent)
	}

	return this
}
