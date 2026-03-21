package responses

type ListAgentsResponse struct {
	PaginateResponse
	Agents []*AgentResponse `json:"agents"`
}
