package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
)

// ExternalListConversationsResponse is the broker's generated list wire shape, named locally to carry mapping.
type ExternalListConversationsResponse clients.ResponsesListConversationsResponse

func (response *ExternalListConversationsResponse) ToUseCaseResponse() *responses.ListBrokerConversationsResponse {
	if response == nil {
		return nil
	}

	var conversations []*responses.BrokerConversationResponse
	if response.Conversations != nil {
		conversations = make([]*responses.BrokerConversationResponse, len(*response.Conversations))
		for i := range *response.Conversations {
			conversation := ExternalConversationResponse((*response.Conversations)[i])
			conversations[i] = conversation.ToUseCaseResponse()
		}
	}

	out := &responses.ListBrokerConversationsResponse{Conversations: conversations}
	if response.Page != nil {
		out.Page = uint(*response.Page)
	}
	if response.ElementsPerPage != nil {
		out.ElementsPerPage = uint(*response.ElementsPerPage)
	}
	if response.TotalElements != nil {
		out.TotalElements = uint(*response.TotalElements)
	}
	return out
}
