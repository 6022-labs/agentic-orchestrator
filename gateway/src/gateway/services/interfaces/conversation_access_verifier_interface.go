package interfaces

import (
	"context"

	"github.com/google/uuid"
)

// ConversationAccessVerifierInterface scopes a conversation to the owners or wallets of its participating agents.
// EnsureAccess fails with NotFoundError (unknown), ForbiddenError (no access) or UnavailableError.
type ConversationAccessVerifierInterface interface {
	EnsureAccess(ctx context.Context, callerId string, conversationId uuid.UUID) error
}
