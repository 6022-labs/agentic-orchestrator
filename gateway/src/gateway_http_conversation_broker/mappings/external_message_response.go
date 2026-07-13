package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/helpers"
	"github.com/samber/lo"
)

// ExternalMessageResponse is the broker's generated message wire shape, named locally to carry mapping.
type ExternalMessageResponse clients.ResponsesConversationMessageResponse

func (message *ExternalMessageResponse) ToUseCaseResponse() *responses.BrokerMessageResponse {
	if message == nil {
		return nil
	}
	return &responses.BrokerMessageResponse{
		ID:             helpers.ParseUuid(message.Id),
		ConversationId: helpers.ParseUuid(message.ConversationId),
		ParticipantId:  helpers.ParseUuid(message.ParticipantId),
		Content:        lo.FromPtr(message.Content),
		CreatedAt:      helpers.ParseTime(message.CreatedAt),
	}
}
