package responses

type AgentComputeResponses struct {
	IsFinal   bool     `json:"isFinal"`
	Responses []string `json:"responses"`
}
