package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
)

// ParticipantRoleResolverInterface reads a participant's role from its on-chain attributes.
type ParticipantRoleResolverInterface interface {
	// Resolve returns nil for users, who have no on-chain identity. Agents always resolve to a role.
	Resolve(ctx context.Context, participant *models.Participant) (*enums.ParticipantRole, error)
}
