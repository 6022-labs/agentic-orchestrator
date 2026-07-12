package responses

import (
	"time"

	"github.com/google/uuid"
)

// BrokerConversationResponse is the broker conversation projection re-exposed to the console.
type BrokerConversationResponse struct {
	ID          uuid.UUID  `json:"id"`
	BridgeId    uuid.UUID  `json:"bridgeId"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	ConcludedAt *time.Time `json:"concludedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
