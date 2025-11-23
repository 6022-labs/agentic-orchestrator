package interfaces

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/google/uuid"
)

type CreateDatabaseConnectionRequest struct {
	AgentId uuid.UUID
	Host    string
	Port    uint16
	Driver  enums.DatabaseConnectionDriver
}

type CreateDatabaseConnectionResult struct {
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
}

type DatabaseConnectionStatus struct {
	IsCompleted bool
	Status      *string
}

type DatabaseConnectionRequesterInterface interface {
	Create(request *CreateDatabaseConnectionRequest) (*CreateDatabaseConnectionResult, error)
	Remove(agentId uuid.UUID) error
	SetupAgentDatabase(agentId uuid.UUID, copyFromAgentId *uuid.UUID) error
	GetSetupStatus(agentId uuid.UUID) (*DatabaseConnectionStatus, error)
}
