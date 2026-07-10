package models

// RemoteAgentCard is the domain view of a node's /.well-known/6022 document.
type RemoteAgentCard struct {
	Protocol        string                   `json:"protocol"`
	ProtocolVersion string                   `json:"protocolVersion"`
	Name            string                   `json:"name"`
	EnsDomain       string                   `json:"ensDomain,omitempty"`
	Endpoints       RemoteAgentCardEndpoints `json:"endpoints"`
	PaymentMethods  []string                 `json:"paymentMethods,omitempty"`
}

type RemoteAgentCardEndpoints struct {
	Responses       string `json:"responses,omitempty"`
	ChatCompletions string `json:"chatCompletions,omitempty"`
}
