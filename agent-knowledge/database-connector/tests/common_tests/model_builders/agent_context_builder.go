package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/google/uuid"
)

type AgentContextBuilder struct {
	id  uuid.UUID
	key string

	createdAt time.Time
	updatedAt time.Time
}

func NewAgentContextBuilder() *AgentContextBuilder {
	id := uuid.New()
	key := uuid.New().String()

	createdAtStr := time.Now().Format(time.RFC3339)
	updatedAtStr := time.Now().Format(time.RFC3339)

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &AgentContextBuilder{
		id:        id,
		key:       key,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (b *AgentContextBuilder) WithID(id uuid.UUID) *AgentContextBuilder {
	b.id = id
	return b
}

func (b *AgentContextBuilder) WithKey(key string) *AgentContextBuilder {
	b.key = key
	return b
}

func (b *AgentContextBuilder) WithCreatedAt(createdAt time.Time) *AgentContextBuilder {
	b.createdAt = createdAt
	return b
}

func (b *AgentContextBuilder) WithUpdatedAt(updatedAt time.Time) *AgentContextBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *AgentContextBuilder) Build() *models.AgentContext {
	return &models.AgentContext{
		ID:        b.id,
		Key:       b.key,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
	}
}
