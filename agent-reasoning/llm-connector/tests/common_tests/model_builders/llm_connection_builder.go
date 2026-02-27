package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type LLMConnectionBuilder struct {
	id          uuid.UUID
	model       string
	baseUrl     *string
	agentId     uuid.UUID
	platform    enums.LLMConnectionPlatform
	setupStatus enums.LLMConnectionSetupStatus
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

func NewLLMConnectionBuilder() *LLMConnectionBuilder {
	id := uuid.New()
	model := faker.Word()

	baseUrl := faker.URL()

	platforms := enums.GetAllLLMConnectionPlatforms()
	platformIndex, _ := faker.RandomInt(0, len(platforms)-1, 1)
	platform := platforms[platformIndex[0]]

	setupStatuses := enums.LLMConnectionSetupStatusEstablished

	agentId := uuid.New()

	createdAtStr := faker.Date()
	updatedAtStr := faker.Date()

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &LLMConnectionBuilder{
		id:          id,
		baseUrl:     &baseUrl,
		model:       model,
		agentId:     agentId,
		platform:    platform,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		setupStatus: setupStatuses,
		deletedAt:   nil,
	}
}

func (b *LLMConnectionBuilder) WithID(id uuid.UUID) *LLMConnectionBuilder {
	b.id = id
	return b
}

func (b *LLMConnectionBuilder) WithBaseUrl(baseUrl string) *LLMConnectionBuilder {
	b.baseUrl = &baseUrl
	return b
}

func (b *LLMConnectionBuilder) WithModel(model string) *LLMConnectionBuilder {
	b.model = model
	return b
}

func (b *LLMConnectionBuilder) WithAgentId(agentId uuid.UUID) *LLMConnectionBuilder {
	b.agentId = agentId
	return b
}

func (b *LLMConnectionBuilder) WithPlatform(platform enums.LLMConnectionPlatform) *LLMConnectionBuilder {
	b.platform = platform
	return b
}

func (b *LLMConnectionBuilder) WithSetupStatus(setupStatus enums.LLMConnectionSetupStatus) *LLMConnectionBuilder {
	b.setupStatus = setupStatus
	return b
}

func (b *LLMConnectionBuilder) WithCreatedAt(createdAt time.Time) *LLMConnectionBuilder {
	b.createdAt = createdAt
	return b
}

func (b *LLMConnectionBuilder) WithUpdatedAt(updatedAt time.Time) *LLMConnectionBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *LLMConnectionBuilder) WithDeletedAt(deletedAt *time.Time) *LLMConnectionBuilder {
	b.deletedAt = deletedAt
	return b
}

func (b *LLMConnectionBuilder) Build() *models.LLMConnection {
	return &models.LLMConnection{
		ID:        b.id,
		Model:     b.model,
		BaseUrl:   b.baseUrl,
		AgentID:   b.agentId,
		Platform:  b.platform,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
		DeletedAt: b.deletedAt,

		SetupStatus: b.setupStatus,
	}
}
