package models

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"

// ExternalWellKnownAgent6022 mirrors a node's /.well-known/6022 document.
type ExternalWellKnownAgent6022 struct {
	Protocol        string                     `json:"protocol"`
	ProtocolVersion string                     `json:"protocolVersion"`
	Name            string                     `json:"name"`
	EnsDomain       string                     `json:"ensDomain,omitempty"`
	Endpoints       ExternalAgent6022Endpoints `json:"endpoints"`
	PaymentMethods  []string                   `json:"paymentMethods,omitempty"`
}

type ExternalAgent6022Endpoints struct {
	Responses       string `json:"responses,omitempty"`
	ChatCompletions string `json:"chatCompletions,omitempty"`
}

func (card *ExternalWellKnownAgent6022) ToModel() *models.RemoteAgentCard {
	if card == nil {
		return nil
	}
	return &models.RemoteAgentCard{
		Protocol:        card.Protocol,
		ProtocolVersion: card.ProtocolVersion,
		Name:            card.Name,
		EnsDomain:       card.EnsDomain,
		Endpoints: models.RemoteAgentCardEndpoints{
			Responses:       card.Endpoints.Responses,
			ChatCompletions: card.Endpoints.ChatCompletions,
		},
		PaymentMethods: card.PaymentMethods,
	}
}
