package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
)

// ExternalListAgentsResponse is the registry's generated list wire shape, named locally to carry mapping.
type ExternalListAgentsResponse clients.ResponsesListAgentsResponse

func (response *ExternalListAgentsResponse) ToUseCaseResponse() *responses.ListAgentsResponse {
	if response == nil {
		return nil
	}

	var agents []*responses.AgentResponse
	if response.Agents != nil {
		agents = make([]*responses.AgentResponse, len(*response.Agents))
		for i := range *response.Agents {
			agent := ExternalAgentResponse((*response.Agents)[i])
			agents[i] = agent.ToUseCaseResponse()
		}
	}

	out := &responses.ListAgentsResponse{Agents: agents}
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
