package responses

// ListAgentsResponse is returned when listing agents
// It contains pagination information and the agents themselves.
type ListAgentsResponse struct {
	PaginateResponse
	Agents []*AgentResponse `json:"agents"`
}
