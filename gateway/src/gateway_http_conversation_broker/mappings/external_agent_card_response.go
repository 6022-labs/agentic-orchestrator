package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/samber/lo"
)

// ExternalAgentCardResponse is the broker's generated 6022-card wire shape, named locally to carry mapping.
type ExternalAgentCardResponse clients.ResponsesRemoteAgentCardResponse

func (card *ExternalAgentCardResponse) ToUseCaseResponse() *responses.AgentCardResponse {
	if card == nil {
		return nil
	}

	endpoints := responses.AgentCardEndpointsResponse{}
	if card.Endpoints != nil {
		endpoints.Responses = lo.FromPtr(card.Endpoints.Responses)
		endpoints.ChatCompletions = lo.FromPtr(card.Endpoints.ChatCompletions)
	}

	var paymentMethods []string
	if card.PaymentMethods != nil {
		paymentMethods = *card.PaymentMethods
	}

	return &responses.AgentCardResponse{
		Protocol:        lo.FromPtr(card.Protocol),
		ProtocolVersion: lo.FromPtr(card.ProtocolVersion),
		Name:            lo.FromPtr(card.Name),
		EnsDomain:       lo.FromPtr(card.EnsDomain),
		Endpoints:       endpoints,
		PaymentMethods:  paymentMethods,
	}
}
