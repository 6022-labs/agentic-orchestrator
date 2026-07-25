package responses

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type ParticipantResponse struct {
	ID                uuid.UUID              `json:"id"`
	CollectionAddress *string                `json:"collectionAddress,omitempty"`
	CollectionTokenId *uint64                `json:"collectionTokenId,omitempty"`
	EnsDomain         *string                `json:"ensDomain,omitempty"`
	Role              *enums.ParticipantRole `json:"role,omitempty"`
	DisplayName       *string                `json:"displayName,omitempty"`
	IsProcessing      bool                   `json:"isProcessing"`
	LastProcessedAt   *time.Time             `json:"lastProcessedAt,omitempty"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

func ParticipantResponseFromModel(participant *models.Participant) *ParticipantResponse {
	if participant == nil {
		return nil
	}
	return &ParticipantResponse{
		ID:                participant.ID,
		CollectionAddress: participant.CollectionAddress,
		CollectionTokenId: participant.CollectionTokenId,
		EnsDomain:         participant.EnsDomain,
		Role:              participant.Role,
		DisplayName:       participant.DisplayName,
		IsProcessing:      participant.IsProcessing,
		LastProcessedAt:   participant.LastProcessedAt,
		CreatedAt:         participant.CreatedAt,
		UpdatedAt:         participant.UpdatedAt,
	}
}
