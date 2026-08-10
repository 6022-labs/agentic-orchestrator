package models

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/google/uuid"
)

// ExternalAgentResponse is the registry's wire shape. We only consume the
// fields we actually need on the broker side.
type ExternalAgentResponse struct {
	ID                uuid.UUID         `json:"id"`
	Name              string            `json:"name"`
	Owner             *string           `json:"owner"`
	EnsDomain         string            `json:"ensDomain"`
	Attributes        map[string]string `json:"attributes"`
	Addresses         map[string]string `json:"addresses"`
	CollectionAddress string            `json:"collectionAddress"`
	CollectionTokenId *uint64           `json:"collectionTokenId"`
}

func (response *ExternalAgentResponse) ToRegistryAgent() *models.RegistryAgent {
	if response == nil {
		return nil
	}
	return &models.RegistryAgent{
		ID:                response.ID,
		Name:              response.Name,
		Owner:             response.Owner,
		EnsDomain:         response.EnsDomain,
		Addresses:         response.Addresses,
		Attributes:        response.Attributes,
		CollectionAddress: response.CollectionAddress,
		CollectionTokenId: response.CollectionTokenId,
	}
}
