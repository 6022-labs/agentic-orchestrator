package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/google/uuid"
)

type AgentContextConversationMappingBuilder struct {
	id             uuid.UUID
	conversationId string

	createdAt time.Time

	agentContextId uuid.UUID
}

func NewAgentContextConversationMappingBuilder() *AgentContextConversationMappingBuilder {
	id := uuid.New()
	conversationId := "test-conversation-id"

	createdAtStr := time.Now().Format(time.RFC3339)

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)

	agentContextId := uuid.New()

	return &AgentContextConversationMappingBuilder{
		id:             id,
		conversationId: conversationId,
		createdAt:      createdAt,
		agentContextId: agentContextId,
	}
}

func (b *AgentContextConversationMappingBuilder) WithID(id uuid.UUID) *AgentContextConversationMappingBuilder {
	b.id = id
	return b
}

func (b *AgentContextConversationMappingBuilder) WithConversationId(conversationId string) *AgentContextConversationMappingBuilder {
	b.conversationId = conversationId
	return b
}

func (b *AgentContextConversationMappingBuilder) WithCreatedAt(createdAt time.Time) *AgentContextConversationMappingBuilder {
	b.createdAt = createdAt
	return b
}

func (b *AgentContextConversationMappingBuilder) WithAgentContextId(agentContextId uuid.UUID) *AgentContextConversationMappingBuilder {
	b.agentContextId = agentContextId
	return b
}

func (b *AgentContextConversationMappingBuilder) Build() *models.AgentContextConversationMapping {
	return &models.AgentContextConversationMapping{
		ID:             b.id,
		CreatedAt:      b.createdAt,
		ConversationId: b.conversationId,
		AgentContextId: b.agentContextId,
	}
}
