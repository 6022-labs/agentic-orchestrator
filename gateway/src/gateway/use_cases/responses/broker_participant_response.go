package responses

import "github.com/google/uuid"

// BrokerParticipantResponse is the broker participant projection re-exposed to the console.
type BrokerParticipantResponse struct {
	ID                uuid.UUID `json:"id"`
	CollectionAddress *string   `json:"collectionAddress,omitempty"`
	CollectionTokenId *uint64   `json:"collectionTokenId,omitempty"`
	EnsDomain         *string   `json:"ensDomain,omitempty"`
	Role              *string   `json:"role,omitempty"`
	DisplayName       *string   `json:"displayName,omitempty"`
}
