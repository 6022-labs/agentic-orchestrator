package responses

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type BridgeResponse struct {
	ID        uuid.UUID        `json:"id"`
	Kind      enums.BridgeKind `json:"kind"`
	BaseUrl   string           `json:"baseUrl"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

func BridgeResponseFromModel(bridge *models.Bridge) *BridgeResponse {
	if bridge == nil {
		return nil
	}
	return &BridgeResponse{
		ID:        bridge.ID,
		Kind:      bridge.Kind,
		BaseUrl:   bridge.BaseUrl,
		CreatedAt: bridge.CreatedAt,
		UpdatedAt: bridge.UpdatedAt,
	}
}
