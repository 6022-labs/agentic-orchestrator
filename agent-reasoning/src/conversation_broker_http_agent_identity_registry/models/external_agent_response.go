package models

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

// ExternalAgentResponse is the registry's wire shape. We only consume the
// fields we actually need on the broker side.
type ExternalAgentResponse struct {
	ID                uuid.UUID         `json:"id"`
	Name              string            `json:"name"`
	Role              string            `json:"role"`
	Owner             *string           `json:"owner"`
	EnsDomain         string            `json:"ensDomain"`
	WalletAddress     string            `json:"walletAddress"`
	Attributes        map[string]string `json:"attributes"`
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
		Role:              response.Role,
		Owner:             response.Owner,
		EnsDomain:         response.EnsDomain,
		WalletAddress:     response.WalletAddress,
		Attributes:        response.Attributes,
		CollectionAddress: response.CollectionAddress,
		CollectionTokenId: response.CollectionTokenId,
	}
}
