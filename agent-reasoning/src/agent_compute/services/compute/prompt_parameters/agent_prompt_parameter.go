package prompt_parameters

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/responses"

type agentPromptParameter struct {
	Name          string            `json:"name"`
	Role          string            `json:"role"`
	Owner         string            `json:"owner"`
	EnsDomain     string            `json:"ensDomain"`
	WalletAddress string            `json:"walletAddress"`
	Attributes    map[string]string `json:"attributes"`
}

func newAgentPromptParameterFromResponse(agent *responses.AgentResponse) *agentPromptParameter {
	owner := ""
	if agent.Owner != nil {
		owner = *agent.Owner
	}
	return &agentPromptParameter{
		Name:          agent.Name,
		Role:          agent.Role.String(),
		Owner:         owner,
		EnsDomain:     agent.EnsDomain,
		Attributes:    agent.Attributes,
		WalletAddress: agent.WalletAddress,
	}
}
