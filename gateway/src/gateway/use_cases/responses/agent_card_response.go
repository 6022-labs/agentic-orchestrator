package responses

// AgentCardResponse mirrors an agent node's /.well-known/6022 card.
type AgentCardResponse struct {
	Protocol        string                     `json:"protocol"`
	ProtocolVersion string                     `json:"protocolVersion"`
	Name            string                     `json:"name"`
	EnsDomain       string                     `json:"ensDomain,omitempty"`
	Endpoints       AgentCardEndpointsResponse `json:"endpoints"`
	PaymentMethods  []string                   `json:"paymentMethods,omitempty"`
}

type AgentCardEndpointsResponse struct {
	Responses       string `json:"responses,omitempty"`
	ChatCompletions string `json:"chatCompletions,omitempty"`
}
