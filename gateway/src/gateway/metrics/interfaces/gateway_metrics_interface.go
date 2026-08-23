package interfaces

import "context"

const (
	AuthMethodJwt       = "jwt"
	AuthMethodSignature = "signature"
	AuthMethodNone      = "none"
)

// GatewayMetricsInterface records the gateway.* and cache.* metrics of the edge.
type GatewayMetricsInterface interface {
	// RecordAuth records one authenticated-route decision; method is none when nothing was injected.
	RecordAuth(ctx context.Context, method string, authenticated bool)
	RecordWalletVerification(ctx context.Context, err error)
	RecordCacheResult(ctx context.Context, cacheName string, hit bool)
}
