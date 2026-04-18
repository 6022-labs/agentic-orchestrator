package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type ConversationMessageBuilder struct {
	id             uuid.UUID
	conversationId uuid.UUID
	participantId  uuid.UUID
	content        string
	createdAt      time.Time
}

func NewConversationMessageBuilder() *ConversationMessageBuilder {
	return &ConversationMessageBuilder{
		id:             uuid.New(),
		conversationId: uuid.New(),
		participantId:  uuid.New(),
		content:        faker.Sentence(),
		createdAt:      time.Now(),
	}
}

func (b *ConversationMessageBuilder) WithID(id uuid.UUID) *ConversationMessageBuilder {
	b.id = id
	return b
}

func (b *ConversationMessageBuilder) WithConversationID(id uuid.UUID) *ConversationMessageBuilder {
	b.conversationId = id
	return b
}

func (b *ConversationMessageBuilder) WithParticipantID(id uuid.UUID) *ConversationMessageBuilder {
	b.participantId = id
	return b
}

func (b *ConversationMessageBuilder) WithContent(s string) *ConversationMessageBuilder {
	b.content = s
	return b
}

func (b *ConversationMessageBuilder) WithCreatedAt(t time.Time) *ConversationMessageBuilder {
	b.createdAt = t
	return b
}

func (b *ConversationMessageBuilder) Build() *models.ConversationMessage {
	return &models.ConversationMessage{
		ID:             b.id,
		ConversationID: b.conversationId,
		ParticipantID:  b.participantId,
		Content:        b.content,
		CreatedAt:      b.createdAt,
	}
}
