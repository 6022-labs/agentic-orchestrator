package prompts

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

type agentPromptParameter struct {
	Name          string            `json:"name"`
	Role          string            `json:"role"`
	Owner         string            `json:"owner"`
	EnsDomain     string            `json:"ensDomain"`
	WalletAddress string            `json:"walletAddress"`
	Attributes    map[string]string `json:"attributes"`
}

func newAgentPromptParameter(agent *models.RegistryAgent) *agentPromptParameter {
	owner := ""
	if agent.Owner != nil {
		owner = *agent.Owner
	}
	return &agentPromptParameter{
		Name:          agent.Name,
		Role:          agent.Role,
		Owner:         owner,
		EnsDomain:     agent.EnsDomain,
		Attributes:    agent.Attributes,
		WalletAddress: agent.WalletAddress,
	}
}

type potentialAgentPromptParameter struct {
	*agentPromptParameter
	Name    string `json:"name"`
	IsAgent bool   `json:"isAgent"`
}

func newPotentialAgentPromptParameter(name string, agent *models.RegistryAgent) *potentialAgentPromptParameter {
	this := &potentialAgentPromptParameter{
		Name:    name,
		IsAgent: agent != nil,
	}

	if agent != nil {
		this.Name = agent.Name
		this.agentPromptParameter = newAgentPromptParameter(agent)
	}

	return this
}
