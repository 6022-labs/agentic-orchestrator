package requests

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
)

type SetupAgentDatabaseSettingsFirstStepRequest struct {
	AgentByCollectionTokenIdRequest
	Host   string                         `form:"host" json:"host"`
	Port   uint16                         `form:"port" json:"port"`
	Driver enums.DatabaseConnectionDriver `form:"driver" json:"driver"`
}

func (r *SetupAgentDatabaseSettingsFirstStepRequest) Validate() error {
	if !enums.GetAllDatabaseConnectionDrivers().Contains(r.Driver) {
		return errors.NewValidationError("driver", "must be one of the valid database drivers: "+enums.GetAllDatabaseConnectionDrivers().String())
	}

	return r.EnsureAuthenticatedWallet()
}
