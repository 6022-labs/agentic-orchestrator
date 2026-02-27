package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type AgentPromptTemplateBuilder struct {
	id      uuid.UUID
	key     enums.AgentPromptTemplateKey
	content string

	createdAt time.Time
	updatedAt time.Time
}

func NewAgentPromptTemplateBuilder() *AgentPromptTemplateBuilder {
	id := uuid.New()
	key := enums.AgentPromptTemplateKey("key")
	content := "content"

	createdAtStr := time.Now().Format(time.RFC3339)
	updatedAtStr := time.Now().Format(time.RFC3339)

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &AgentPromptTemplateBuilder{
		id:        id,
		key:       key,
		content:   content,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (b *AgentPromptTemplateBuilder) WithID(id uuid.UUID) *AgentPromptTemplateBuilder {
	b.id = id
	return b
}

func (b *AgentPromptTemplateBuilder) WithKey(key enums.AgentPromptTemplateKey) *AgentPromptTemplateBuilder {
	b.key = key
	return b
}

func (b *AgentPromptTemplateBuilder) WithContent(content string) *AgentPromptTemplateBuilder {
	b.content = content
	return b
}

func (b *AgentPromptTemplateBuilder) WithCreatedAt(createdAt time.Time) *AgentPromptTemplateBuilder {
	b.createdAt = createdAt
	return b
}

func (b *AgentPromptTemplateBuilder) WithUpdatedAt(updatedAt time.Time) *AgentPromptTemplateBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *AgentPromptTemplateBuilder) Build() *models.AgentPromptTemplate {
	return &models.AgentPromptTemplate{
		ID:        b.id,
		Key:       b.key,
		Content:   b.content,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
	}
}
