package models

import "github.com/google/uuid"

type ExternalRegisterAgentRequest struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionTokenId string `json:"collectionTokenId"`
}

type ExternalRegisterUserRequest struct {
	DisplayName string `json:"displayName"`
}

type ExternalParticipantResponse struct {
	ID                uuid.UUID `json:"id"`
	CollectionAddress *string   `json:"collectionAddress,omitempty"`
	CollectionTokenId *uint64   `json:"collectionTokenId,omitempty"`
	EnsDomain         *string   `json:"ensDomain,omitempty"`
	Role              *string   `json:"role,omitempty"`
	DisplayName       *string   `json:"displayName,omitempty"`
}
