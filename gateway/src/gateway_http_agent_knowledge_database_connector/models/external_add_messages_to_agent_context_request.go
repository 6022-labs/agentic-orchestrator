package models

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

type ExternalMessage struct {
	Author         string                            `json:"author"`
	Content        string                            `json:"content"`
	Platform       enums.AgentContextMessagePlatform `json:"platform"`
	Role           enums.AgentContextMessageRole     `json:"role"`
	ConversationId *string                           `json:"conversationId" optional:"true"`
}

type ExternalAddMessagesToAgentContextRequest struct {
	AgentId        uuid.UUID         `json:"agentId"`
	ContextKey     *string           `json:"contextKey,omitempty"`
	ConversationId *string           `json:"conversationId,omitempty"`
	Messages       []ExternalMessage `json:"messages"`
}
