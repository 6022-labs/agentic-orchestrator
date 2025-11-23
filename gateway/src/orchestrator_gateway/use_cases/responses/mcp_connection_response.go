package responses

type MCPConnectionResponse struct {
	ID      string `json:"id"`
	AgentID string `json:"agentId"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}
