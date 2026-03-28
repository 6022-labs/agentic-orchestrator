package models

import "github.com/google/uuid"

type GetLinkedAgentContextOfConversationResponseStatus string

const (
	GetLinkedAgentContextOfConversationResponseStatusNotFound GetLinkedAgentContextOfConversationResponseStatus = "not-found"
	GetLinkedAgentContextOfConversationResponseStatusFound    GetLinkedAgentContextOfConversationResponseStatus = "found"
)

type ExternalGetLinkedAgentContextOfConversationResponse struct {
	Id     uuid.UUID                                         `json:"id"`
	Key    string                                            `json:"key"`
	Status GetLinkedAgentContextOfConversationResponseStatus `json:"status"`
}
