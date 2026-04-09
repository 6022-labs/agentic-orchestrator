package requests

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
)

type GetAgentPromptTemplateRequest struct {
	AgentByCollectionTokenIdRequest
	Key enums.AgentPromptTemplateKey `json:"key" params:"key"`
}

func (r *GetAgentPromptTemplateRequest) Validate() error {
	if !enums.GetAllAgentPromptTemplateKeys().Contains(r.Key) {
		return errors.NewValidationError("key", "key is invalid, must be one of: "+enums.GetAllAgentPromptTemplateKeys().String())
	}

	return r.EnsureAuthenticatedWallet()
}
