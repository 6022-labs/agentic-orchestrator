package models

import "github.com/google/uuid"

// ExternalAgentResponse is the registry's wire shape. We only consume the
// fields we actually need on the broker side.
type ExternalAgentResponse struct {
	ID                uuid.UUID `json:"id"`
	Role              string    `json:"role"`
	EnsDomain         string    `json:"ensDomain"`
	CollectionAddress string    `json:"collectionAddress"`
	CollectionTokenId *uint64   `json:"collectionTokenId"`
}
