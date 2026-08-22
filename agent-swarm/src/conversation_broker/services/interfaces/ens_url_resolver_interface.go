package interfaces

import "context"

// EnsUrlResolverInterface resolves an agent ENS name to its runtime base URL by reading the
// AgentEnsRegistry on the single chain this broker is configured for.
type EnsUrlResolverInterface interface {
	// ResolveUrl errors on a name belonging to another chain, and returns "" when no url is published.
	ResolveUrl(ctx context.Context, ensDomain string) (string, error)
}
