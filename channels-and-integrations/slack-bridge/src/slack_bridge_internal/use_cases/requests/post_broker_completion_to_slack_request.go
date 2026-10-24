package requests

import (
	"strings"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type PostBrokerCompletionToSlackRequest struct {
	ConversationId    uuid.UUID `json:"conversationId"`
	CollectionAddress string    `json:"collectionAddress"`
	CollectionTokenId uint64    `json:"collectionTokenId"`
	Content           string    `json:"content"`
}

// ValidateAndSanitize lowercases the collection address so it matches the stored rows.
func (r *PostBrokerCompletionToSlackRequest) ValidateAndSanitize() error {
	if r.ConversationId == uuid.Nil {
		return common_errors.NewValidationError("conversationId", "is required")
	}
	if !strings.HasPrefix(r.CollectionAddress, "0x") || !common.IsHexAddress(r.CollectionAddress) {
		return common_errors.NewValidationError("collectionAddress", "must be a 0x-prefixed hex address")
	}
	if strings.TrimSpace(r.Content) == "" {
		return common_errors.NewValidationError("content", "is required")
	}

	r.CollectionAddress = strings.ToLower(r.CollectionAddress)
	return nil
}
