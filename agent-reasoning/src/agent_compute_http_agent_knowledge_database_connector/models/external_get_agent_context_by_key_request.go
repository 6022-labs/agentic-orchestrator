package models

import "github.com/google/uuid"

type ExternalGetAgentContextByKeyRequest struct {
	AgentId uuid.UUID
	Key     string
}
