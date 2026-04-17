package requests

import (
	"errors"
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
		return errors.New("conversationId is required")
	}
	if strings.TrimSpace(r.CollectionAddress) == "" {
		return errors.New("collectionAddress is required")
	}
	if strings.TrimSpace(r.Content) == "" {
		return errors.New("content is required")
	}
	return nil
}
