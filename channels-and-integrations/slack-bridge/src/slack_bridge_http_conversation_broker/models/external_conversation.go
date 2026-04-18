package models

import "github.com/google/uuid"

type ExternalRegisterConversationRequest struct {
	BridgeID       uuid.UUID   `json:"bridgeId"`
	Description    *string     `json:"description,omitempty"`
	ParticipantIDs []uuid.UUID `json:"participantIds,omitempty"`
}

type ExternalUpdateConversationDescriptionRequest struct {
	Description *string `json:"description"`
}

type ExternalConversationParticipantRequest struct {
	ParticipantID uuid.UUID `json:"participantId"`
}

type ExternalIngestMessageRequest struct {
	ParticipantID uuid.UUID `json:"participantId"`
	Content       string    `json:"content"`
}

type ExternalConversationResponse struct {
	ID          uuid.UUID `json:"id"`
	BridgeID    uuid.UUID `json:"bridgeId"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
}

type ExternalConversationMessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversationId"`
	ParticipantID  uuid.UUID `json:"participantId"`
	Content        string    `json:"content"`
}
