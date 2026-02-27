package requests

import (
	"strings"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/errors"
	"github.com/ethereum/go-ethereum/common"
)

type UpsertAgentImageRequest struct {
	Key               string `json:"key"`
	NewValue          string `json:"newValue"`
	CollectionTokenId uint64 `json:"collectionTokenId"`
	CollectionAddress string `json:"collectionAddress"`
}

func (r *UpsertAgentImageRequest) Validate() error {
	if r.CollectionTokenId == 0 {
		return errors.NewValidationError("collectionTokenId", "is required")
	}

	if !common.IsHexAddress(r.CollectionAddress) {
		return errors.NewValidationError("collectionAddress", "is not a valid hex address")
	}

	if len(strings.TrimSpace(r.Key)) == 0 {
		return errors.NewValidationError("key", "is required")
	}

	if len(strings.TrimSpace(r.NewValue)) == 0 {
		return errors.NewValidationError("newValue", "is required")
	}

	return nil
}
