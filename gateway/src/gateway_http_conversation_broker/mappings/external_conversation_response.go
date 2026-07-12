package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/helpers"
)

// ExternalConversationResponse is the broker's generated conversation wire shape, named locally to carry mapping.
type ExternalConversationResponse clients.ResponsesConversationResponse

func (conversation *ExternalConversationResponse) ToUseCaseResponse() *responses.BrokerConversationResponse {
	if conversation == nil {
		return nil
	}

	var status string
	if conversation.Status != nil {
		status = string(*conversation.Status)
	}

	return &responses.BrokerConversationResponse{
		ID:          helpers.ParseUuid(conversation.Id),
		BridgeId:    helpers.ParseUuid(conversation.BridgeId),
		Description: conversation.Description,
		Status:      status,
		ConcludedAt: helpers.ParseTimePtr(conversation.ConcludedAt),
		CreatedAt:   helpers.ParseTime(conversation.CreatedAt),
		UpdatedAt:   helpers.ParseTime(conversation.UpdatedAt),
	}
}
