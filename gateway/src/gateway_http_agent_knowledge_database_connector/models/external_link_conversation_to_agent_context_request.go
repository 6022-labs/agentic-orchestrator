package models

import "github.com/google/uuid"

type ExternalLinkConversationToAgentContextRequest struct {
	Key            string    `json:"key"`
	AgentId        uuid.UUID `json:"agentId"`
	ConversationId string    `json:"conversationId"`
}
