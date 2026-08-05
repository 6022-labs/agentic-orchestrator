package responses

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/google/uuid"
)

type AgentResponse struct {
	ID uuid.UUID `json:"id"`

	Name       string                       `json:"name"`
	Owner      *string                      `json:"owner"`
	Creator    string                       `json:"creator"`
	Images     map[string]string            `json:"images"`
	CloneOf    *uint64                      `json:"cloneOf"`
	Attributes map[string]string            `json:"attributes"`
	EnsDomain  string                       `json:"ensDomain"`
	Addresses  map[enums.AddressType]string `json:"addresses"`

	ProposalId        *uint64 `json:"proposalId"`
	CollectionTokenId *uint64 `json:"collectionTokenId"`
	CollectionAddress string  `json:"collectionAddress"`
}

func AgentResponseFromModel(agent *models.Agent, agentAddresses []*models.AgentAddress) *AgentResponse {
	if agent == nil {
		return nil
	}

	addresses := make(map[enums.AddressType]string, len(agentAddresses))
	for _, agentAddress := range agentAddresses {
		addresses[agentAddress.Type] = agentAddress.Address
	}

	return &AgentResponse{
		ID:                agent.ID,
		Name:              agent.Name,
		Owner:             agent.Owner,
		Images:            agent.Images,
		Creator:           agent.Creator,
		CloneOf:           agent.CloneOf,
		EnsDomain:         agent.EnsDomain,
		Attributes:        agent.Attributes,
		Addresses:         addresses,
		ProposalId:        agent.ProposalId,
		CollectionTokenId: agent.CollectionTokenId,
		CollectionAddress: agent.CollectionAddress,
	}
}
