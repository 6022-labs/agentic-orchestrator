package metrics

import (
	"context"
	"time"

	common_metrics "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const meterName = "agentic_orchestrator_conversation_broker_remote_agent"

// providerName is required by semconv; the remote node picks its own model vendor, so this names the dialect.
const providerName = "6022.agent"

type RemoteAgentMetrics struct {
	invocationHistogram metric.Float64Histogram
	paymentSettlements  metric.Int64Counter
}

func NewRemoteAgentMetrics() *RemoteAgentMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	invocationHistogram, err := meter.Float64Histogram(
		"gen_ai.client.operation.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.GenAiDurationSecondsBuckets...),
		metric.WithDescription("Duration of GenAI client operations"),
	)
	if err != nil {
		invocationHistogram, _ = noopMeter.Float64Histogram("gen_ai.client.operation.duration")
	}

	paymentSettlements, err := meter.Int64Counter(
		"payment.settlement",
		metric.WithDescription("x402 payment settlements attempted, by outcome"),
	)
	if err != nil {
		paymentSettlements, _ = noopMeter.Int64Counter("payment.settlement")
	}

	return &RemoteAgentMetrics{
		invocationHistogram: invocationHistogram,
		paymentSettlements:  paymentSettlements,
	}
}

func (m *RemoteAgentMetrics) RecordInvocation(
	ctx context.Context,
	agentName string,
	endpointUrl string,
	duration time.Duration,
	err error,
) {
	attrs := []attribute.KeyValue{
		attribute.String("gen_ai.operation.name", "invoke_agent"),
		attribute.String("gen_ai.provider.name", providerName),
		attribute.String("gen_ai.agent.name", agentName),
	}
	attrs = append(attrs, common_metrics.ServerAttributes(endpointUrl)...)
	// error.type is conditionally required, and only on a failed operation.
	if err != nil {
		attrs = append(attrs, attribute.String("error.type", common_metrics.ErrorType(err)))
	}
	m.invocationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *RemoteAgentMetrics) RecordPaymentSettlement(ctx context.Context, outcome string, network string) {
	m.paymentSettlements.Add(ctx, 1, metric.WithAttributes(
		attribute.String("outcome", outcome),
		attribute.String("payment.network", network),
		// This side pays; agent-node records the receiving side.
		attribute.String("payment.direction", "outbound"),
	))
}
