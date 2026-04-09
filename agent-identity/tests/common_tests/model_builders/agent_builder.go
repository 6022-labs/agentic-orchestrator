package model_builders

import (
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type AgentBuilder struct {
	id                uuid.UUID
	name              string
	role              enums.AgentRole
	owner             *string
	creator           string
	walletAddress     string
	isInternal        bool
	images            map[string]string
	attributes        map[string]string
	ensDomain         string
	collectionAddress string
	collectionTokenId *uint64
	proposalId        *uint64
	createdAt         time.Time
	updatedAt         time.Time
	deletedAt         *time.Time
}

func NewAgentBuilder() *AgentBuilder {
	id := uuid.New()
	name := faker.Name()
	ensDomain := fmt.Sprintf("%s.6022.eth", name)
	role := enums.AgentRoleExpert

	ids, _ := faker.RandomInt(1, 100, 2)
	collectionTokenId := uint64(ids[0])
	proposalId := uint64(ids[1])

	privateKey, _ := crypto.GenerateKey()
	owner := crypto.PubkeyToAddress(privateKey.PublicKey).String()
	creator := owner

	createdAtStr := faker.Date()
	updatedAtStr := faker.Date()

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &AgentBuilder{
		id:                id,
		name:              name,
		role:              role,
		owner:             &owner,
		creator:           creator,
		ensDomain:         ensDomain,
		isInternal:        false,
		images: map[string]string{
			"default": "QmTzQ1Z5b7g8d9f8e6c3b2a1c4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z",
		},
		attributes:        map[string]string{},
		updatedAt:         updatedAt,
		createdAt:         createdAt,
		walletAddress:     common_tests.RandomEthAddress().Hex(),
		proposalId:        &proposalId,
		collectionTokenId: &collectionTokenId,
		collectionAddress: common_tests.RandomEthAddress().Hex(),
		deletedAt:         nil,
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

func (b *AgentBuilder) WithRole(role enums.AgentRole) *AgentBuilder {
	b.role = role
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

func (b *AgentBuilder) WithIsInternal(isInternal bool) *AgentBuilder {
	b.isInternal = isInternal
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

func (b *AgentBuilder) WithWalletAddress(walletAddress string) *AgentBuilder {
	b.walletAddress = walletAddress
	return b
}

func (b *AgentBuilder) WithCreatedAt(createdAt time.Time) *AgentBuilder {
	b.createdAt = createdAt
	return b
}

func (b *AgentBuilder) WithUpdatedAt(updatedAt time.Time) *AgentBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *AgentBuilder) WithDeletedAt(deletedAt *time.Time) *AgentBuilder {
	b.deletedAt = deletedAt
	return b
}

func (b *AgentBuilder) Build() *models.Agent {
	return &models.Agent{
		ID:                b.id,
		Name:              b.name,
		Role:              b.role,
		Owner:             b.owner,
		Images:            b.images,
		Creator:           b.creator,
		EnsDomain:         b.ensDomain,
		IsInternal:        b.isInternal,
		ProposalId:        b.proposalId,
		Attributes:        b.attributes,
		WalletAddress:     b.walletAddress,
		CollectionTokenId: b.collectionTokenId,
		CreatedAt:         b.createdAt,
		UpdatedAt:         b.updatedAt,
		DeletedAt:         b.deletedAt,
	}
}
