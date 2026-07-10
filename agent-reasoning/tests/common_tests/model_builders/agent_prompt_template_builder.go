package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type AgentPromptTemplateBuilder struct {
	id                uuid.UUID
	collectionAddress string
	collectionTokenId uint64
	key               enums.AgentPromptTemplateKey
	content           string
	createdAt         time.Time
	updatedAt         time.Time
}

func NewAgentPromptTemplateBuilder() *AgentPromptTemplateBuilder {
	return &AgentPromptTemplateBuilder{
		id:                uuid.New(),
		collectionAddress: "0x" + faker.UUIDDigit(),
		collectionTokenId: 1,
		key:               enums.AgentPromptTemplateKeyCloneSystem,
		content:           faker.Sentence(),
		createdAt:         time.Now(),
		updatedAt:         time.Now(),
	}
}

func (b *AgentPromptTemplateBuilder) WithCollectionAddress(address string) *AgentPromptTemplateBuilder {
	b.collectionAddress = address
	return b
}

func (b *AgentPromptTemplateBuilder) WithCollectionTokenId(tokenId uint64) *AgentPromptTemplateBuilder {
	b.collectionTokenId = tokenId
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

func (b *AgentPromptTemplateBuilder) Build() *models.AgentPromptTemplate {
	return &models.AgentPromptTemplate{
		ID:                b.id,
		CollectionAddress: b.collectionAddress,
		CollectionTokenId: b.collectionTokenId,
		Key:               b.key,
		Content:           b.content,
		CreatedAt:         b.createdAt,
		UpdatedAt:         b.updatedAt,
	}
}
