package responses

import (
	"time"

	"github.com/google/uuid"
)

// BrokerMessageResponse is the broker message projection re-exposed to the console.
type BrokerMessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationId uuid.UUID `json:"conversationId"`
	ParticipantId  uuid.UUID `json:"participantId"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}
