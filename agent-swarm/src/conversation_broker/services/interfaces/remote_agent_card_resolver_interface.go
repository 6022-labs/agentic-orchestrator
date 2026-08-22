package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

// RemoteAgentCardResolverInterface resolves an agent's ENS domain to its
// published 6022 card, caching results with a configurable TTL.
type RemoteAgentCardResolverInterface interface {
	ResolveCard(ctx context.Context, ensDomain string) (*models.RemoteAgentCard, error)
}
