package interfaces

import "github.com/google/uuid"

type ConversationBrokerRequesterInterface interface {
	RegisterAgent(collectionAddress string, collectionTokenId uint64) (uuid.UUID, error)
	RegisterUser(displayName string) (uuid.UUID, error)
	RegisterConversation(bridgeId uuid.UUID, description *string, participantIds []uuid.UUID) (uuid.UUID, error)
	IngestMessage(conversationId uuid.UUID, participantId uuid.UUID, content string) (uuid.UUID, error)
	ForceConcludeConversation(conversationId uuid.UUID) error
	AttachParticipant(conversationId uuid.UUID, participantId uuid.UUID) error
	DetachParticipant(conversationId uuid.UUID, participantId uuid.UUID) error
}
