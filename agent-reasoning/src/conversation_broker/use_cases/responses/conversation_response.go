package responses

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type ConversationResponse struct {
	ID            uuid.UUID                `json:"id"`
	BridgeID      uuid.UUID                `json:"bridgeId"`
	Description   *string                  `json:"description,omitempty"`
	Status        enums.ConversationStatus `json:"status"`
	ConcludedAt   *time.Time               `json:"concludedAt,omitempty"`
	UnconcludedAt *time.Time               `json:"unconcludedAt,omitempty"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
}

func ConversationResponseFromModel(conversation *models.Conversation) *ConversationResponse {
	if conversation == nil {
		return nil
	}
	return &ConversationResponse{
		ID:            conversation.ID,
		BridgeID:      conversation.BridgeID,
		Description:   conversation.Description,
		Status:        conversation.Status,
		ConcludedAt:   conversation.ConcludedAt,
		UnconcludedAt: conversation.UnconcludedAt,
		CreatedAt:     conversation.CreatedAt,
		UpdatedAt:     conversation.UpdatedAt,
	}
}
