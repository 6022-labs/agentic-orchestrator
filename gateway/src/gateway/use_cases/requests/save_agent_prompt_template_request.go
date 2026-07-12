package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
)

type SaveAgentPromptTemplateRequest struct {
	GetAgentPromptTemplateRequest
	Content string `json:"content"`
}

func (r *SaveAgentPromptTemplateRequest) Validate() error {
	if err := r.GetAgentPromptTemplateRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(r.Content) == "" {
		return errors.NewValidationError("content", "must not be empty")
	}
	return nil
}
