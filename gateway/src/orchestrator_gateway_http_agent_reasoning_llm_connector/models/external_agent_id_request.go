package models

import "github.com/google/uuid"

type ExternalAgentIdRequest struct {
	AgentId uuid.UUID `json:"agentId"`
}
