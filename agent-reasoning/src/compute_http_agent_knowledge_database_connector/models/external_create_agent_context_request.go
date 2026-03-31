package models

import "github.com/google/uuid"

type ExternalCreateAgentContextRequest struct {
	AgentId    uuid.UUID `json:"agentId"`
	Key        string    `json:"key"`
	IsInternal bool      `json:"isInternal"`
}
