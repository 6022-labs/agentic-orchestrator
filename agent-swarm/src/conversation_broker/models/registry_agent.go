package models

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const (
	// EvmAddressType keys the agent's Ethereum wallet inside Addresses.
	EvmAddressType = "evm"
)

// RegistryAgent is the registry's agent projection consumed by the broker
// (participant registration and prompt-parameter resolution).
type RegistryAgent struct {
	ID                uuid.UUID         `json:"id"`
	Name              string            `json:"name"`
	Owner             *string           `json:"owner"`
	EnsDomain         string            `json:"ensDomain"`
	Addresses         map[string]string `json:"addresses"`
	Attributes        map[string]string `json:"attributes"`
	CollectionAddress string            `json:"collectionAddress"`
	CollectionTokenId *uint64           `json:"collectionTokenId"`
}

// EvmAddress is the key that signs the agent's card and settles its payments.
func (agent *RegistryAgent) EvmAddress() string {
	return agent.Addresses[EvmAddressType]
}

func (agent *RegistryAgent) ToParticipant() *Participant {
	if agent == nil {
		return nil
	}

	return &Participant{
		ID:                agent.ID,
		CollectionTokenId: agent.CollectionTokenId,
		EnsDomain:         lo.ToPtr(agent.EnsDomain),
		CollectionAddress: lo.ToPtr(agent.CollectionAddress),
	}
}
