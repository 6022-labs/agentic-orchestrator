package requests

import (
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
)

type RegisterMintedParticipantRequest struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionTokenId string `json:"collectionTokenId"`
}

func (r *RegisterMintedParticipantRequest) Validate() error {
	if strings.TrimSpace(r.CollectionAddress) == "" {
		return errors.NewValidationError("collectionAddress", "must not be empty")
	}
	if strings.TrimSpace(r.CollectionTokenId) == "" {
		return errors.NewValidationError("collectionTokenId", "must not be empty")
	}
	if _, err := strconv.ParseUint(r.CollectionTokenId, 10, 64); err != nil {
		return errors.NewValidationError("collectionTokenId", "must be a valid unsigned integer")
	}
	return nil
}

func (r *RegisterMintedParticipantRequest) ParsedCollectionTokenId() uint64 {
	v, _ := strconv.ParseUint(r.CollectionTokenId, 10, 64)
	return v
}
