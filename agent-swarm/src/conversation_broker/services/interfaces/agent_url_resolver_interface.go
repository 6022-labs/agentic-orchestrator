package interfaces

import "context"

// AgentUrlResolverInterface resolves an agent's node URL from its ENS url
// record, caching the result per domain.
type AgentUrlResolverInterface interface {
	ResolveUrl(ctx context.Context, ensDomain string) (string, error)
}
