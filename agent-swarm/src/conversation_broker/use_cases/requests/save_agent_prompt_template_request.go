package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
)

type SaveAgentPromptTemplateRequest struct {
	GetAgentPromptTemplateRequest
	Content string `json:"content"`
}

func (r *SaveAgentPromptTemplateRequest) ValidateAndSanitize() error {
	if err := r.GetAgentPromptTemplateRequest.ValidateAndSanitize(); err != nil {
		return err
	}
	if strings.TrimSpace(r.Content) == "" {
		return errors.NewValidationError("content", "must not be empty")
	}
	return nil
}
