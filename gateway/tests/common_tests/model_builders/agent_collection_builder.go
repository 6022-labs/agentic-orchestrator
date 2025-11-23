package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/tests/common_tests"
	"github.com/go-faker/faker/v4"
)

type AgentCollectionBuilder struct {
	address   string
	name      string
	symbol    string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewAgentCollectionBuilder() *AgentCollectionBuilder {
	name := faker.Name()
	symbol := faker.Word()
	createdAt := time.Now()
	updatedAt := time.Now()
	var deletedAt *time.Time

	return &AgentCollectionBuilder{
		address:   common_tests.RandomEthAddress().Hex(),
		name:      name,
		symbol:    symbol,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (b *AgentCollectionBuilder) WithAddress(address string) *AgentCollectionBuilder {
	b.address = address
	return b
}

func (b *AgentCollectionBuilder) WithName(name string) *AgentCollectionBuilder {
	b.name = name
	return b
}

func (b *AgentCollectionBuilder) WithSymbol(symbol string) *AgentCollectionBuilder {
	b.symbol = symbol
	return b
}

func (b *AgentCollectionBuilder) WithCreatedAt(createdAt time.Time) *AgentCollectionBuilder {
	b.createdAt = createdAt
	return b
}

func (b *AgentCollectionBuilder) WithUpdatedAt(updatedAt time.Time) *AgentCollectionBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *AgentCollectionBuilder) WithDeletedAt(deletedAt *time.Time) *AgentCollectionBuilder {
	b.deletedAt = deletedAt
	return b
}

func (b *AgentCollectionBuilder) Build() *models.AgentCollection {
	return &models.AgentCollection{
		Address:   b.address,
		Name:      b.name,
		Symbol:    b.symbol,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
		DeletedAt: b.deletedAt,
	}
}
