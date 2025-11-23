package models

import (
	"github.com/google/uuid"
)

type ExternalChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ExternalChatCompletionRequest struct {
	AgentId  uuid.UUID                        `json:"agentId"`
	Messages []*ExternalChatCompletionMessage `json:"messages"`
}
