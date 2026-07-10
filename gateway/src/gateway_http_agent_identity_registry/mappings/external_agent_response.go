package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_identity_registry/helpers"
)

// ExternalAgentResponse is the registry's generated agent wire shape, named locally to carry mapping.
type ExternalAgentResponse clients.ResponsesAgentResponse

func (response *ExternalAgentResponse) ToUseCaseResponse() *responses.AgentResponse {
	if response == nil {
		return nil
	}

	var role enums.AgentRole
	if response.Role != nil {
		role = enums.AgentRole(*response.Role)
	}

	return &responses.AgentResponse{
		ID:                helpers.ParseUuid(response.Id),
		Name:              helpers.DerefString(response.Name),
		Role:              role,
		Owner:             response.Owner,
		Creator:           helpers.DerefString(response.Creator),
		Images:            helpers.DerefStringMap(response.Images),
		CloneOf:           helpers.IntToUint64Ptr(response.CloneOf),
		Attributes:        helpers.DerefStringMap(response.Attributes),
		EnsDomain:         helpers.DerefString(response.EnsDomain),
		WalletAddress:     helpers.DerefString(response.WalletAddress),
		ProposalId:        helpers.IntToUint64Ptr(response.ProposalId),
		CollectionTokenId: helpers.IntToUint64Ptr(response.CollectionTokenId),
		CollectionAddress: helpers.DerefString(response.CollectionAddress),
	}
}
