package prompt_parameters

import "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"

type agentPromptParameter struct {
	Name          string            `json:"name"`
	Role          string            `json:"role"`
	Owner         string            `json:"owner"`
	EnsDomain     string            `json:"ensDomain"`
	WalletAddress string            `json:"walletAddress"`
	Attributes    map[string]string `json:"attributes"`
}

func newAgentPromptParameterFromModel(agent *models.Agent) *agentPromptParameter {
	return &agentPromptParameter{
		Name:          agent.Name,
		Role:          agent.Role.String(),
		Owner:         *agent.Owner,
		EnsDomain:     agent.EnsDomain,
		Attributes:    agent.Attributes,
		WalletAddress: agent.WalletAddress,
	}
}
