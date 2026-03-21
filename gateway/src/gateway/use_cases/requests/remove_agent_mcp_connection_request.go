package requests

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/google/uuid"
)

type RemoveAgentMCPConnectionRequest struct {
	AgentByCollectionTokenIdRequest
	MCPConnectionId uuid.UUID `json:"mcpConnectionId" params:"mcpConnectionId"`
}

func (r *RemoveAgentMCPConnectionRequest) Validate() error {
	if r.MCPConnectionId == uuid.Nil {
		return errors.NewValidationError("mcpConnectionId", "MCP connection ID is required")
	}

	return nil
}
