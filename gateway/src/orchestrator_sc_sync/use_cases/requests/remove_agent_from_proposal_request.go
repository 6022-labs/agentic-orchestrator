package requests

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/ethereum/go-ethereum/common"
)

type RemoveAgentFromProposalRequest struct {
	ProposalId        uint64 `json:"proposalId"`
	CollectionAddress string `json:"collectionAddress"`
}

func (r *RemoveAgentFromProposalRequest) Validate() error {
	if r.ProposalId == 0 {
		return errors.NewValidationError("proposalId", "is required")
	}

	if !common.IsHexAddress(r.CollectionAddress) {
		return errors.NewValidationError("collectionAddress", "is not a valid hex address")
	}

	return nil
}
