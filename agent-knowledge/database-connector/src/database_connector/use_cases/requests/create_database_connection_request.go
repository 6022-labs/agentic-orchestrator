package requests

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type CreateDatabaseConnectionRequest struct {
	Host    string                         `json:"host"`
	Port    uint16                         `json:"port"`
	Driver  enums.DatabaseConnectionDriver `json:"driver"`
	AgentID uuid.UUID                      `json:"agentId"`
}

func (request *CreateDatabaseConnectionRequest) Validate() error {
	if len(strings.TrimSpace(request.Host)) == 0 {
		return errors.NewValidationError("host", "host is required")
	}

	if !enums.GetAllDatabaseConnectionDrivers().Contains(request.Driver) {
		return errors.NewValidationError("driver", fmt.Sprintf("driver %s is required and must be one of: %s", request.Driver, enums.GetAllDatabaseConnectionDrivers().String()))
	}

	if request.AgentID == uuid.Nil {
		return errors.NewValidationError("agentId", "agentId is required")
	}

	return nil
}
