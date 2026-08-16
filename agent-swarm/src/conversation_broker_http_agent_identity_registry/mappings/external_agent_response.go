package mappings

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_http_agent_identity_registry/helpers"
	"github.com/samber/lo"
)

// ExternalAgentResponse is the registry's generated agent wire shape, named locally to carry mapping.
type ExternalAgentResponse clients.ResponsesAgentResponse

func (response *ExternalAgentResponse) ToRegistryAgent() *models.RegistryAgent {
	if response == nil {
		return nil
	}

	return &models.RegistryAgent{
		ID:                helpers.ParseUuid(response.Id),
		Name:              lo.FromPtr(response.Name),
		Owner:             response.Owner,
		EnsDomain:         lo.FromPtr(response.EnsDomain),
		Addresses:         lo.FromPtr(response.Addresses),
		Attributes:        lo.FromPtr(response.Attributes),
		CollectionAddress: lo.FromPtr(response.CollectionAddress),
		CollectionTokenId: helpers.IntToUint64Ptr(response.CollectionTokenId),
	}
}
