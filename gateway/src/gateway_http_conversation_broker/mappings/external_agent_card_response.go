package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/helpers"
)

// ExternalAgentCardResponse is the broker's generated 6022-card wire shape, named locally to carry mapping.
type ExternalAgentCardResponse clients.ModelsRemoteAgentCard

func (card *ExternalAgentCardResponse) ToUseCaseResponse() *responses.AgentCardResponse {
	if card == nil {
		return nil
	}

	endpoints := responses.AgentCardEndpointsResponse{}
	if card.Endpoints != nil {
		endpoints.Responses = helpers.DerefString(card.Endpoints.Responses)
		endpoints.ChatCompletions = helpers.DerefString(card.Endpoints.ChatCompletions)
	}

	var paymentMethods []string
	if card.PaymentMethods != nil {
		paymentMethods = *card.PaymentMethods
	}

	return &responses.AgentCardResponse{
		Protocol:        helpers.DerefString(card.Protocol),
		ProtocolVersion: helpers.DerefString(card.ProtocolVersion),
		Name:            helpers.DerefString(card.Name),
		EnsDomain:       helpers.DerefString(card.EnsDomain),
		Endpoints:       endpoints,
		PaymentMethods:  paymentMethods,
	}
}
