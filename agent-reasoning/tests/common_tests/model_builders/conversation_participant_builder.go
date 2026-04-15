package model_builders

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationParticipantBuilder struct {
	conversationId uuid.UUID
	participantId  uuid.UUID
	joinedAt       time.Time
}

func NewConversationParticipantBuilder() *ConversationParticipantBuilder {
	return &ConversationParticipantBuilder{
		conversationId: uuid.New(),
		participantId:  uuid.New(),
		joinedAt:       time.Now(),
	}
}

func (b *ConversationParticipantBuilder) WithConversationID(id uuid.UUID) *ConversationParticipantBuilder {
	b.conversationId = id
	return b
}

func (b *ConversationParticipantBuilder) WithParticipantID(id uuid.UUID) *ConversationParticipantBuilder {
	b.participantId = id
	return b
}

func (b *ConversationParticipantBuilder) Build() *models.ConversationParticipant {
	return &models.ConversationParticipant{
		ConversationID: b.conversationId,
		ParticipantID:  b.participantId,
		JoinedAt:       b.joinedAt,
	}
}
