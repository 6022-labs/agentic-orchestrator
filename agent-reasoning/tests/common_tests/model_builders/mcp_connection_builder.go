package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type MCPConnectionBuilder struct {
	id        uuid.UUID
	name      string
	url       string
	agentId   uuid.UUID
	createdAt time.Time
	updatedAt time.Time

	headerAuthKeys []string
}

func NewMCPConnectionBuilder() *MCPConnectionBuilder {
	id := uuid.New()
	name := faker.Word()

	url := faker.URL()

	headerAuthKeys := []string{"authorization"}
	agentId := uuid.New()

	createdAtStr := faker.Date()
	updatedAtStr := faker.Date()

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &MCPConnectionBuilder{
		id:             id,
		url:            url,
		name:           name,
		agentId:        agentId,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		headerAuthKeys: headerAuthKeys,
	}
}

func (b *MCPConnectionBuilder) WithID(id uuid.UUID) *MCPConnectionBuilder {
	b.id = id
	return b
}

func (b *MCPConnectionBuilder) WithName(name string) *MCPConnectionBuilder {
	b.name = name
	return b
}

func (b *MCPConnectionBuilder) WithUrl(url string) *MCPConnectionBuilder {
	b.url = url
	return b
}

func (b *MCPConnectionBuilder) WithHeaderAuthKeys(headerAuthKeys []string) *MCPConnectionBuilder {
	b.headerAuthKeys = headerAuthKeys
	return b
}

func (b *MCPConnectionBuilder) WithAgentId(agentId uuid.UUID) *MCPConnectionBuilder {
	b.agentId = agentId
	return b
}

func (b *MCPConnectionBuilder) WithCreatedAt(createdAt time.Time) *MCPConnectionBuilder {
	b.createdAt = createdAt
	return b
}

func (b *MCPConnectionBuilder) WithUpdatedAt(updatedAt time.Time) *MCPConnectionBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *MCPConnectionBuilder) Build() *models.MCPConnection {
	// Convert headerAuthKeys to MCPConnectionHeaderAuthKeys
	var headerAuthKeys models.MCPConnectionHeaderAuthKeys
	headerAuthKeys.FromArray(b.headerAuthKeys)

	return &models.MCPConnection{
		ID:             b.id,
		Url:            b.url,
		Name:           b.name,
		AgentID:        b.agentId,
		CreatedAt:      b.createdAt,
		UpdatedAt:      b.updatedAt,
		HeaderAuthKeys: headerAuthKeys,
	}
}
