package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type ConversationBrokerRequesterInterface interface {
	RegisterAgent(ctx context.Context, collectionAddress string, collectionTokenId uint64) (uuid.UUID, error)
	RegisterUser(ctx context.Context, displayName string) (uuid.UUID, error)
	RegisterConversation(
		ctx context.Context,
		bridgeId uuid.UUID,
		description *string,
		participantIds []uuid.UUID,
	) (uuid.UUID, error)
	IngestMessage(
		ctx context.Context,
		conversationId uuid.UUID,
		participantId uuid.UUID,
		content string,
	) (uuid.UUID, error)
	ForceConcludeConversation(ctx context.Context, conversationId uuid.UUID) error
	AttachParticipant(ctx context.Context, conversationId uuid.UUID, participantId uuid.UUID) error
	DetachParticipant(ctx context.Context, conversationId uuid.UUID, participantId uuid.UUID) error
}
