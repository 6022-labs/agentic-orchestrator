package models

import "github.com/google/uuid"

// ExternalCompletionCallbackRequest is the payload broker POSTs to a bridge's
// base_url + "/broker/callback" whenever an agent has produced a completion
// in a conversation. Bridges resolve `ConversationId` via their own local
// link table and identify the agent via its NFT identity (collection_address
// + collection_token_id) — no broker_participant_id mapping required.
type ExternalCompletionCallbackRequest struct {
	ConversationId    uuid.UUID `json:"conversationId"`
	CollectionAddress string    `json:"collectionAddress"`
	CollectionTokenId uint64    `json:"collectionTokenId"`
	Content           string    `json:"content"`
}
