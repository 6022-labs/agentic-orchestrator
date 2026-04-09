package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
)

type AddMessageToAgentContextRequest struct {
	AgentByCollectionTokenIdRequest
	ContextKey string                        `params:"contextKey" swaggerignore:"true"`
	Message    string                        `json:"message"`
	Role       enums.AgentContextMessageRole `json:"role"`
}

func (r *AddMessageToAgentContextRequest) Validate() error {
	if len(strings.TrimSpace(r.ContextKey)) < 2 {
		return errors.NewValidationError("contextKey", "Context key need to be at least 3 characters long")
	}

	if len(strings.TrimSpace(r.Message)) < 2 {
		return errors.NewValidationError("message", "Message need to be at least 3 characters long")
	}

	if !enums.GetAllAgentContextMessageRoles().Contains(r.Role) {
		return errors.NewValidationError("role", "Role must be one of: "+enums.GetAllAgentContextMessageRoles().String())
	}

	return r.EnsureAuthenticatedWallet()
}
