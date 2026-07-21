package model_builders

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type AgentBuilder struct {
	id                uuid.UUID
	name              string
	owner             *string
	creator           string
	addresses         map[string]string
	images            map[string]string
	attributes        map[string]string
	ensDomain         string
	collectionAddress string
	collectionTokenId *uint64
	proposalId        *uint64
}

func NewAgentBuilder() *AgentBuilder {
	id := uuid.New()
	name := faker.Name()
	ensDomain := fmt.Sprintf("%s.6022.eth", name)

	ids, _ := faker.RandomInt(1, 100, 2)
	collectionTokenId := uint64(ids[0])
	proposalId := uint64(ids[1])

	privateKey, _ := crypto.GenerateKey()
	owner := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	creator := owner

	return &AgentBuilder{
		id:        id,
		name:      name,
		owner:     &owner,
		creator:   creator,
		ensDomain: ensDomain,
		images: map[string]string{
			"default": "QmTzQ1Z5b7g8d9f8e6c3b2a1c4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z",
		},
		attributes:        map[string]string{"role": "expert"},
		addresses:         map[string]string{"evm": common_tests.RandomEthAddress().Hex()},
		proposalId:        &proposalId,
		collectionTokenId: &collectionTokenId,
		collectionAddress: common_tests.RandomEthAddress().Hex(),
	}
}

func (b *AgentBuilder) WithID(id uuid.UUID) *AgentBuilder {
	b.id = id
	return b
}

func (b *AgentBuilder) WithName(name string) *AgentBuilder {
	b.name = name
	return b
}

func (b *AgentBuilder) WithEnsDomain(ensDomain string) *AgentBuilder {
	b.ensDomain = ensDomain
	return b
}

func (b *AgentBuilder) WithOwner(owner string) *AgentBuilder {
	b.owner = &owner
	return b
}

func (b *AgentBuilder) WithCreator(creator string) *AgentBuilder {
	b.creator = creator
	return b
}

func (b *AgentBuilder) WithImages(images map[string]string) *AgentBuilder {
	b.images = images
	return b
}

func (b *AgentBuilder) WithAttributes(attributes map[string]string) *AgentBuilder {
	b.attributes = attributes
	return b
}

func (b *AgentBuilder) WithProposalId(proposalId *uint64) *AgentBuilder {
	b.proposalId = proposalId
	return b
}

func (b *AgentBuilder) WithCollectionTokenId(collectionTokenId *uint64) *AgentBuilder {
	b.collectionTokenId = collectionTokenId
	return b
}

func (b *AgentBuilder) WithCollectionAddress(collectionAddress string) *AgentBuilder {
	b.collectionAddress = collectionAddress
	return b
}

func (b *AgentBuilder) WithAddresses(addresses map[string]string) *AgentBuilder {
	b.addresses = addresses
	return b
}

func (b *AgentBuilder) Build() *responses.AgentResponse {
	return &responses.AgentResponse{
		ID:                b.id,
		Name:              b.name,
		Owner:             b.owner,
		Images:            b.images,
		Creator:           b.creator,
		EnsDomain:         b.ensDomain,
		ProposalId:        b.proposalId,
		Attributes:        b.attributes,
		Addresses:         b.addresses,
		CollectionTokenId: b.collectionTokenId,
		CollectionAddress: b.collectionAddress,
	}
}
