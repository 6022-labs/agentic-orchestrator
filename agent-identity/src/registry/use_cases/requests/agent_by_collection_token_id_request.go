package requests

import (
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
)

type AgentByCollectionTokenIdRequest struct {
	AgentCollectionAddress string `param:"agentCollectionAddress"`
	AgentCollectionTokenId uint64 `param:"agentCollectionTokenId"`
}

// ValidateAndSanitize lowercases the collection address so it matches the stored rows.
func (r *AgentByCollectionTokenIdRequest) ValidateAndSanitize() error {
	if !strings.HasPrefix(r.AgentCollectionAddress, "0x") || !common.IsHexAddress(r.AgentCollectionAddress) {
		return common_errors.NewValidationError("agentCollectionAddress", "must be a 0x-prefixed hex address")
	}
	// Token IDs start at 1 on-chain, so 0 can never identify an agent.
	if r.AgentCollectionTokenId == 0 {
		return common_errors.NewValidationError("agentCollectionTokenId", "must be a positive integer")
	}

	r.AgentCollectionAddress = strings.ToLower(r.AgentCollectionAddress)
	return nil
}
