package responses

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
)

type RemoteAgentCardEndpointsResponse struct {
	Responses       string `json:"responses,omitempty"`
	ChatCompletions string `json:"chatCompletions,omitempty"`
}

type RemoteAgentCardResponse struct {
	Protocol        string                           `json:"protocol"`
	ProtocolVersion string                           `json:"protocolVersion"`
	Name            string                           `json:"name"`
	EnsDomain       string                           `json:"ensDomain,omitempty"`
	Endpoints       RemoteAgentCardEndpointsResponse `json:"endpoints"`
	PaymentMethods  []string                         `json:"paymentMethods,omitempty"`
}

func RemoteAgentCardResponseFromModel(card *models.RemoteAgentCard) *RemoteAgentCardResponse {
	if card == nil {
		return nil
	}
	return &RemoteAgentCardResponse{
		Protocol:        card.Protocol,
		ProtocolVersion: card.ProtocolVersion,
		Name:            card.Name,
		EnsDomain:       card.EnsDomain,
		Endpoints: RemoteAgentCardEndpointsResponse{
			Responses:       card.Endpoints.Responses,
			ChatCompletions: card.Endpoints.ChatCompletions,
		},
		PaymentMethods: card.PaymentMethods,
	}
}
