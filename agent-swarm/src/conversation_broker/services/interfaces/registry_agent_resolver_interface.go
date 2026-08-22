package interfaces

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
)

// RegistryAgentResolverInterface caches the registry's agent projection.
// Returns nil (no error) when the agent is not found.
type RegistryAgentResolverInterface interface {
	ResolveByCollection(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
	) (*models.RegistryAgent, error)
	ResolveByEnsDomain(ctx context.Context, ensDomain string) (*models.RegistryAgent, error)
}
