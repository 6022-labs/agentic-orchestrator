package models

import "github.com/google/uuid"

// ExternalCompletionCallbackRequest is the payload broker POSTs to a bridge's
// callback_url whenever an agent has produced a completion in a conversation.
// Bridges resolve `ConversationId` and `ParticipantId` via their own local
// linking tables to determine where/how to deliver the message.
type ExternalCompletionCallbackRequest struct {
	ConversationId uuid.UUID `json:"conversationId"`
	ParticipantId  uuid.UUID `json:"participantId"`
	Content        string    `json:"content"`
}
