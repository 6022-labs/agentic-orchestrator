package prompts

type agentPromptParameter struct {
	Name       string            `json:"name"`
	Owner      string            `json:"owner"`
	EnsDomain  string            `json:"ensDomain"`
	Addresses  map[string]string `json:"addresses"`
	Attributes map[string]string `json:"attributes"`
}

type potentialAgentPromptParameter struct {
	*agentPromptParameter
	Name    string `json:"name"`
	IsAgent bool   `json:"isAgent"`
}
