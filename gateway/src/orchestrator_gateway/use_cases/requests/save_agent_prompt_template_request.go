package requests

import (
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
)

type SaveAgentPromptTemplateRequest struct {
	AgentByCollectionTokenIdRequest
	Key     enums.AgentPromptTemplateKey `json:"key" params:"key"`
	Content string                       `json:"content"`
}

func (r *SaveAgentPromptTemplateRequest) Validate() error {
	if len(strings.TrimSpace(r.Content)) == 0 {
		return errors.NewValidationError("content", "content is required")
	}

	if !enums.GetAllAgentPromptTemplateKeys().Contains(r.Key) {
		return errors.NewValidationError("key", "key is invalid, must be one of: "+enums.GetAllAgentPromptTemplateKeys().String())
	}

	return r.EnsureAuthenticatedWallet()
}
