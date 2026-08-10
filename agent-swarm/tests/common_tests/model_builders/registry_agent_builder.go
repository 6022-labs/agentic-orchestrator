package model_builders

import (
	"fmt"
	"math/rand"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type RegistryAgentBuilder struct {
	id                uuid.UUID
	name              string
	owner             *string
	ensDomain         string
	addresses         map[string]string
	attributes        map[string]string
	collectionAddress string
	collectionTokenId *uint64
}

func NewRegistryAgentBuilder() *RegistryAgentBuilder {
	collectionTokenId := rand.Uint64()

	return &RegistryAgentBuilder{
		id:                uuid.New(),
		name:              faker.Word(),
		ensDomain:         faker.Word() + ".eth",
		addresses:         map[string]string{models.EvmAddressType: fmt.Sprintf("0x%040x", rand.Uint64())},
		attributes:        map[string]string{"role": enums.ParticipantRoleClone.String()},
		collectionAddress: fmt.Sprintf("0x%040x", rand.Uint64()),
		collectionTokenId: &collectionTokenId,
	}
}

func (b *RegistryAgentBuilder) WithID(id uuid.UUID) *RegistryAgentBuilder {
	b.id = id
	return b
}

func (b *RegistryAgentBuilder) WithName(name string) *RegistryAgentBuilder {
	b.name = name
	return b
}

func (b *RegistryAgentBuilder) WithOwner(owner string) *RegistryAgentBuilder {
	b.owner = &owner
	return b
}

func (b *RegistryAgentBuilder) WithEnsDomain(ensDomain string) *RegistryAgentBuilder {
	b.ensDomain = ensDomain
	return b
}

// WithRole sets the role attribute, which is where the registry actually publishes it.
func (b *RegistryAgentBuilder) WithRole(role string) *RegistryAgentBuilder {
	b.attributes["role"] = role
	return b
}

func (b *RegistryAgentBuilder) WithoutRole() *RegistryAgentBuilder {
	delete(b.attributes, "role")
	return b
}

func (b *RegistryAgentBuilder) WithAttribute(key, value string) *RegistryAgentBuilder {
	b.attributes[key] = value
	return b
}

func (b *RegistryAgentBuilder) WithAddress(addressType, address string) *RegistryAgentBuilder {
	b.addresses[addressType] = address
	return b
}

func (b *RegistryAgentBuilder) WithCollection(address string, tokenId uint64) *RegistryAgentBuilder {
	b.collectionAddress = address
	b.collectionTokenId = &tokenId
	return b
}

func (b *RegistryAgentBuilder) Build() *models.RegistryAgent {
	return &models.RegistryAgent{
		ID:                b.id,
		Name:              b.name,
		Owner:             b.owner,
		EnsDomain:         b.ensDomain,
		Addresses:         b.addresses,
		Attributes:        b.attributes,
		CollectionAddress: b.collectionAddress,
		CollectionTokenId: b.collectionTokenId,
	}
}
