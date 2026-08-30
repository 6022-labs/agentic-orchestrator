package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

// MemorySessionResolverInterface answers which node-side session a turn runs on.
type MemorySessionResolverInterface interface {
	// Resolve records a derived session; nil means an expert has not named one yet.
	Resolve(
		ctx context.Context,
		role *enums.ParticipantRole,
		conversation *models.Conversation,
		participant *models.Participant,
		participants []*models.Participant,
	) (*uuid.UUID, error)
	// Bind records the session an expert named in its reply.
	Bind(
		ctx context.Context,
		conversation *models.Conversation,
		participant *models.Participant,
		sessionId uuid.UUID,
	) error
}
