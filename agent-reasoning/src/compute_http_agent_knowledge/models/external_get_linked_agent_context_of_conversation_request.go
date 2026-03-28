package models

import "github.com/google/uuid"

type ExternalGetLinkedAgentContextOfConversationRequest struct {
	AgentId        uuid.UUID `json:"agentId" query:"agentId"`
	ConversationId string    `json:"conversationId" query:"conversationId"`
}
