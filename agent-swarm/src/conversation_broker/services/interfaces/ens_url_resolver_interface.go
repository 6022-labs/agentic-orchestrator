package interfaces

// EnsUrlResolverInterface resolves an agent ENS name to its runtime
// base URL via ERC-3668 CCIP-Read against the offchain resolver.
type EnsUrlResolverInterface interface {
	ResolveUrl(ensDomain string) (string, error)
}
