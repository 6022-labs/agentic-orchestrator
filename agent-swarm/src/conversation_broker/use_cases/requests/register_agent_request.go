package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
)

type RegisterAgentRequest struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionTokenId uint64 `json:"collectionTokenId,string"`
}

// ValidateAndSanitize lowercases the collection address so it matches the stored rows.
func (r *RegisterAgentRequest) ValidateAndSanitize() error {
	r.CollectionAddress = strings.TrimSpace(r.CollectionAddress)
	if !strings.HasPrefix(r.CollectionAddress, "0x") || !common.IsHexAddress(r.CollectionAddress) {
		return errors.NewValidationError("collectionAddress", "must be a 0x-prefixed hex address")
	}
	// Token IDs start at 1 on-chain, so 0 can never identify an agent.
	if r.CollectionTokenId == 0 {
		return errors.NewValidationError("collectionTokenId", "must be a positive integer")
	}

	r.CollectionAddress = strings.ToLower(r.CollectionAddress)
	return nil
}
