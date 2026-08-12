package interfaces

// EnsUrlResolverInterface resolves an agent ENS name to its runtime base URL by reading the
// AgentEnsRegistry on the single chain this broker is configured for.
type EnsUrlResolverInterface interface {
	// ResolveUrl errors on a name belonging to another chain, and returns "" when no url is published.
	ResolveUrl(ensDomain string) (string, error)
}
