package interfaces

import (
	"context"

	"github.com/google/uuid"
)

// ConversationConcludeInterface marks a conversation as concluded.
//   - Conclude: returns NotFoundError when the conversation is missing.
//   - ConcludeIfExists: best-effort variant; missing conversation is a no-op.
//
// Both are idempotent: an already-concluded conversation is left untouched
// (preserves the original concluded_at).
type ConversationConcludeInterface interface {
	Conclude(ctx context.Context, conversationId uuid.UUID) error
	ConcludeIfExists(ctx context.Context, conversationId uuid.UUID) error
}
