package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/samber/lo"
)

// ExternalAgentCollectionResponse is the registry's generated collection wire shape, named locally to carry mapping.
type ExternalAgentCollectionResponse clients.ResponsesAgentCollectionResponse

func (response *ExternalAgentCollectionResponse) ToUseCaseResponse() *responses.AgentCollectionResponse {
	if response == nil {
		return nil
	}
	return &responses.AgentCollectionResponse{
		Address: lo.FromPtr(response.Address),
		Name:    lo.FromPtr(response.Name),
		Symbol:  lo.FromPtr(response.Symbol),
	}
}
