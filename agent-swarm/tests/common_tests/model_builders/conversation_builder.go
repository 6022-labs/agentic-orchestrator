package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type ConversationBuilder struct {
	id          uuid.UUID
	bridgeId    uuid.UUID
	description *string
	status      enums.ConversationStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewConversationBuilder() *ConversationBuilder {
	description := faker.Sentence()
	return &ConversationBuilder{
		id:          uuid.New(),
		bridgeId:    uuid.New(),
		description: &description,
		status:      enums.ConversationStatusActive,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (b *ConversationBuilder) WithID(id uuid.UUID) *ConversationBuilder {
	b.id = id
	return b
}

func (b *ConversationBuilder) WithBridgeID(id uuid.UUID) *ConversationBuilder {
	b.bridgeId = id
	return b
}

func (b *ConversationBuilder) WithStatus(s enums.ConversationStatus) *ConversationBuilder {
	b.status = s
	return b
}

func (b *ConversationBuilder) Build() *models.Conversation {
	return &models.Conversation{
		ID:          b.id,
		BridgeID:    b.bridgeId,
		Description: b.description,
		Status:      b.status,
		CreatedAt:   b.createdAt,
		UpdatedAt:   b.updatedAt,
	}
}
