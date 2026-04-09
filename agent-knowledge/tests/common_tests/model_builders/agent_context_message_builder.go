package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type AgentContextMessageBuilder struct {
	id      uuid.UUID
	content string
	role    enums.AgentContextMessageRole

	createdAt time.Time

	agentContextId uuid.UUID
}

func NewAgentContextMessageBuilder() *AgentContextMessageBuilder {
	id := uuid.New()
	content := "message"
	role := enums.AgentContextMessageRoleUser

	createdAtStr := time.Now().Format(time.RFC3339)

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)

	agentContextId := uuid.New()

	return &AgentContextMessageBuilder{
		id:             id,
		role:           role,
		content:        content,
		createdAt:      createdAt,
		agentContextId: agentContextId,
	}
}

func (b *AgentContextMessageBuilder) WithID(id uuid.UUID) *AgentContextMessageBuilder {
	b.id = id
	return b
}

func (b *AgentContextMessageBuilder) WithContent(content string) *AgentContextMessageBuilder {
	b.content = content
	return b
}

func (b *AgentContextMessageBuilder) WithRole(role enums.AgentContextMessageRole) *AgentContextMessageBuilder {
	b.role = role
	return b
}

func (b *AgentContextMessageBuilder) WithCreatedAt(createdAt time.Time) *AgentContextMessageBuilder {
	b.createdAt = createdAt
	return b
}

func (b *AgentContextMessageBuilder) WithAgentContextId(agentContextId uuid.UUID) *AgentContextMessageBuilder {
	b.agentContextId = agentContextId
	return b
}

func (b *AgentContextMessageBuilder) Build() *models.AgentContextMessage {
	return &models.AgentContextMessage{
		ID:             b.id,
		Role:           b.role,
		Content:        b.content,
		CreatedAt:      b.createdAt,
		AgentContextId: b.agentContextId,
	}
}
