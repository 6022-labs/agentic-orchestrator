package models

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/google/uuid"
)

type ExternalCreateDatabaseConnectionRequest struct {
	AgentId uuid.UUID                      `json:"agentId"`
	Driver  enums.DatabaseConnectionDriver `json:"driver"`
	Host    string                         `json:"host"`
	Port    uint16                         `json:"port"`
}
