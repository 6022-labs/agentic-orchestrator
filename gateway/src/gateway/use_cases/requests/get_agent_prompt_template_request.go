package requests

import (
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/helpers"
)

type GetAgentPromptTemplateRequest struct {
	AgentCollectionAddress string `uri:"agentCollectionAddress"`
	AgentCollectionTokenId string `uri:"agentCollectionTokenId"`
	Key                    string `uri:"key"`
}

func (r *GetAgentPromptTemplateRequest) ValidateAndSanitize() error {
	r.AgentCollectionAddress = strings.TrimSpace(r.AgentCollectionAddress)
	if !helpers.IsEthAddress(r.AgentCollectionAddress) {
		return errors.NewValidationError("agentCollectionAddress", "must be a 0x-prefixed 20-byte hex address")
	}
	r.AgentCollectionAddress = strings.ToLower(r.AgentCollectionAddress)
	// Token IDs start at 1 on-chain, so 0 can never identify an agent.
	if tokenId, err := strconv.ParseUint(r.AgentCollectionTokenId, 10, 64); err != nil || tokenId == 0 {
		return errors.NewValidationError("agentCollectionTokenId", "must be a positive integer")
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
