package requests

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PostBrokerCompletionToSlackRequest struct {
	ConversationId    uuid.UUID `json:"conversationId"`
	CollectionAddress string    `json:"collectionAddress"`
	CollectionTokenId uint64    `json:"collectionTokenId"`
	Content           string    `json:"content"`
}

func (r *PostBrokerCompletionToSlackRequest) Validate() error {
	if r.ConversationId == uuid.Nil {
		return fmt.Errorf("%w: conversationId is required", ErrValidation)
	}
	if strings.TrimSpace(r.CollectionAddress) == "" {
		return fmt.Errorf("%w: collectionAddress is required", ErrValidation)
	}
	if strings.TrimSpace(r.Content) == "" {
		return fmt.Errorf("%w: content is required", ErrValidation)
	}
	return nil
}
