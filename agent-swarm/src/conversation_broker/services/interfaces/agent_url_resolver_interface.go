package interfaces

// AgentUrlResolverInterface resolves an agent's node URL from its ENS url
// record, caching the result per domain.
type AgentUrlResolverInterface interface {
	ResolveUrl(ensDomain string) (string, error)
}
