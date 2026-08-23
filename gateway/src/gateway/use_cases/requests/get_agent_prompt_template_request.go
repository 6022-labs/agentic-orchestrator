package requests

import (
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
)

type GetAgentPromptTemplateRequest struct {
	AgentCollectionAddress string `uri:"agentCollectionAddress"`
	AgentCollectionTokenId string `uri:"agentCollectionTokenId"`
	Key                    string `uri:"key"`
}

func (r *GetAgentPromptTemplateRequest) Validate() error {
	if strings.TrimSpace(r.AgentCollectionAddress) == "" {
		return errors.NewValidationError("agentCollectionAddress", "must not be empty")
	}
	if _, err := strconv.ParseUint(r.AgentCollectionTokenId, 10, 64); err != nil {
		return errors.NewValidationError("agentCollectionTokenId", "must be a valid unsigned integer")
	}
	if strings.TrimSpace(r.Key) == "" {
		return errors.NewValidationError("key", "must not be empty")
	}
	return nil
}

func (r *GetAgentPromptTemplateRequest) ParsedCollectionTokenId() uint64 {
	v, _ := strconv.ParseUint(r.AgentCollectionTokenId, 10, 64)
	return v
}
