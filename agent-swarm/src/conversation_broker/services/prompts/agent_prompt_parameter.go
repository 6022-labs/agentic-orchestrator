package prompts

type agentPromptParameter struct {
	Name          string            `json:"name"`
	Role          string            `json:"role"`
	Owner         string            `json:"owner"`
	EnsDomain     string            `json:"ensDomain"`
	WalletAddress string            `json:"walletAddress"`
	Attributes    map[string]string `json:"attributes"`
}

type potentialAgentPromptParameter struct {
	*agentPromptParameter
	Name    string `json:"name"`
	IsAgent bool   `json:"isAgent"`
}
