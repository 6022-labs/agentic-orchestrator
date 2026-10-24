package requests

import (
	"fmt"
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
		return fmt.Errorf("%w: conversationId is required", common_errors.ErrValidation)
	}
	if !strings.HasPrefix(r.CollectionAddress, "0x") || !common.IsHexAddress(r.CollectionAddress) {
		return fmt.Errorf("%w: collectionAddress must be a 0x-prefixed hex address", common_errors.ErrValidation)
	}
	if strings.TrimSpace(r.Content) == "" {
		return fmt.Errorf("%w: content is required", common_errors.ErrValidation)
	}

	r.CollectionAddress = strings.ToLower(r.CollectionAddress)
	return nil
}
