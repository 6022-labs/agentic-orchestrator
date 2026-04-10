package model_builders

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/use_cases/responses"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type AgentResponseBuilder struct {
	id                uuid.UUID
	name              string
	role              enums.AgentRole
	owner             *string
	creator           string
	ensDomain         string
	walletAddress     string
	isInternal        bool
	images            map[string]string
	attributes        map[string]string
	collectionAddress string
	collectionTokenId *uint64
	proposalId        *uint64
}

func NewAgentResponseBuilder() *AgentResponseBuilder {
	id := uuid.New()
	name := faker.Name()
	owner := faker.Word()
	creator := faker.Word()
	ensDomain := faker.DomainName()
	walletAddress := faker.Word()
	collectionAddress := faker.Word()
	collectionTokenId := uint64(1)

	return &AgentResponseBuilder{
		id:                id,
		name:              name,
		role:              enums.AgentRoleExpert,
		owner:             &owner,
		creator:           creator,
		ensDomain:         ensDomain,
		walletAddress:     walletAddress,
		isInternal:        false,
		images:            map[string]string{},
		attributes:        map[string]string{},
		collectionAddress: collectionAddress,
		collectionTokenId: &collectionTokenId,
		proposalId:        nil,
	}
}

func (b *AgentResponseBuilder) WithRole(role enums.AgentRole) *AgentResponseBuilder {
	b.role = role
	return b
}

func (b *AgentResponseBuilder) WithEnsDomain(ensDomain string) *AgentResponseBuilder {
	b.ensDomain = ensDomain
	return b
}

func (b *AgentResponseBuilder) Build() *responses.AgentResponse {
	return &responses.AgentResponse{
		ID:                b.id,
		Name:              b.name,
		Role:              b.role,
		Owner:             b.owner,
		Creator:           b.creator,
		EnsDomain:         b.ensDomain,
		WalletAddress:     b.walletAddress,
		IsInternal:        b.isInternal,
		Images:            b.images,
		Attributes:        b.attributes,
		CollectionAddress: b.collectionAddress,
		CollectionTokenId: b.collectionTokenId,
		ProposalId:        b.proposalId,
	}
}
