package responses

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationMessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversationId"`
	ParticipantID  uuid.UUID `json:"participantId"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ConversationMessageResponseFromModel(message *models.ConversationMessage) *ConversationMessageResponse {
	if message == nil {
		return nil
	}
	return &ConversationMessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		ParticipantID:  message.ParticipantID,
		Content:        message.Content,
		CreatedAt:      message.CreatedAt,
	}
}
