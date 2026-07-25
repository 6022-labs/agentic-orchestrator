package models

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// RegistryAgent is the registry's agent projection consumed by the broker
// (participant registration and prompt-parameter resolution).
type RegistryAgent struct {
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

func (agent *RegistryAgent) ToParticipant() (*Participant, error) {
	if agent == nil {
		return nil, nil
	}

	var role enums.ParticipantRole
	if err := role.Set(agent.Role); err != nil {
		return nil, fmt.Errorf("registry returned invalid role: %w", err)
	}

	return &Participant{
		Role:              &role,
		ID:                agent.ID,
		CollectionTokenId: agent.CollectionTokenId,
		EnsDomain:         lo.ToPtr(agent.EnsDomain),
		CollectionAddress: lo.ToPtr(agent.CollectionAddress),
	}, nil
}
