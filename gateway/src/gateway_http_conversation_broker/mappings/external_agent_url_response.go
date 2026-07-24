package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_conversation_broker/clients"
	"github.com/samber/lo"
)

// ExternalAgentUrlResponse is the broker's generated agent-url wire shape, named locally to carry mapping.
type ExternalAgentUrlResponse clients.ResponsesAgentUrlResponse

func (url *ExternalAgentUrlResponse) ToUseCaseResponse() *responses.AgentUrlResponse {
	if url == nil {
		return nil
	}

	return &responses.AgentUrlResponse{
		Url: lo.FromPtr(url.Url),
	}
}
