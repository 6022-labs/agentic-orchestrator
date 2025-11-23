package requests

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/ethereum/go-ethereum/common"
)

type PairProposalToAgentRequest struct {
	ProposalId        uint64 `json:"proposalId"`
	CollectionAddress string `json:"collectionAddress"`
}

func (r *PairProposalToAgentRequest) Validate() error {
	if r.ProposalId == 0 {
		return errors.NewValidationError("proposalId", "is required")
	}

	if !common.IsHexAddress(r.CollectionAddress) {
		return errors.NewValidationError("collectionAddress", "is not a valid hex address")
	}

	return nil
}
