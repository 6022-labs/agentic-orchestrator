package models

import "github.com/google/uuid"

type ExternalGetAgentContextMessagesRequest struct {
	AgentId uuid.UUID `json:"agentId"`
	Key     string    `json:"key"`
}
