package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const meterName = "agentic_orchestrator_gateway"

// errorTypeOther is the semconv catch-all for an error with no lower-cardinality classification.
const errorTypeOther = "_OTHER"

type GatewayMetrics struct {
	authSuccessCounter      metric.Int64Counter
	authFailureCounter      metric.Int64Counter
	walletVerificationCount metric.Int64Counter
	cacheResultCounter      metric.Int64Counter
}

func NewGatewayMetrics() *GatewayMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	authSuccessCounter, err := meter.Int64Counter(
		"gateway.auth.success",
		metric.WithDescription("Total number of requests that reached an authenticated route"),
	)
	if err != nil {
		authSuccessCounter, _ = noopMeter.Int64Counter("gateway.auth.success")
	}

	authFailureCounter, err := meter.Int64Counter(
		"gateway.auth.failure",
		metric.WithDescription("Total number of requests rejected on an authenticated route"),
	)
	if err != nil {
		authFailureCounter, _ = noopMeter.Int64Counter("gateway.auth.failure")
	}

	walletVerificationCount, err := meter.Int64Counter(
		"gateway.wallet.verification",
		metric.WithDescription("Wallet ownership verifications for a bridge, by outcome"),
	)
	if err != nil {
		walletVerificationCount, _ = noopMeter.Int64Counter("gateway.wallet.verification")
	}

	cacheResultCounter, err := meter.Int64Counter(
		"cache.result",
		metric.WithDescription("Cache lookups by hit/miss"),
	)
	if err != nil {
		cacheResultCounter, _ = noopMeter.Int64Counter("cache.result")
	}

	return &GatewayMetrics{
		authSuccessCounter:      authSuccessCounter,
		authFailureCounter:      authFailureCounter,
		walletVerificationCount: walletVerificationCount,
		cacheResultCounter:      cacheResultCounter,
	}
}

func (m *GatewayMetrics) RecordAuth(ctx context.Context, method string, authenticated bool) {
	attrs := metric.WithAttributes(attribute.String("auth.method", method))
	if authenticated {
		m.authSuccessCounter.Add(ctx, 1, attrs)
		return
	}
	m.authFailureCounter.Add(ctx, 1, attrs)
}

func (m *GatewayMetrics) RecordWalletVerification(ctx context.Context, err error) {
	attrs := []attribute.KeyValue{}
	// error.type is conditionally required, and only on a failed operation.
	if err != nil {
		attrs = append(attrs, attribute.String("error.type", errorTypeOther))
	}
	m.walletVerificationCount.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *GatewayMetrics) RecordCacheResult(ctx context.Context, cacheName string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	m.cacheResultCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache.name", cacheName),
		attribute.String("result", result),
	))
}
