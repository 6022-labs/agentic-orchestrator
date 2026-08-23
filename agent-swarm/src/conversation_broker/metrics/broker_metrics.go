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

const meterName = "agentic_orchestrator_conversation_broker"

// errorTypeOther is the semconv catch-all for an error with no lower-cardinality classification.
const errorTypeOther = "_OTHER"

type BrokerMetrics struct {
	conversationsStarted   metric.Int64Counter
	conversationsConcluded metric.Int64Counter
	messagesIngested       metric.Int64Counter
	turnHistogram          metric.Float64Histogram
	activeTurns            metric.Int64UpDownCounter
	daemonTicks            metric.Int64Counter
	agentResolutions       metric.Int64Counter
	cacheResults           metric.Int64Counter
}

func NewBrokerMetrics() *BrokerMetrics {
	noopMeter := noop.NewMeterProvider().Meter(meterName)
	meter := otel.GetMeterProvider().Meter(meterName)

	conversationsStarted, err := meter.Int64Counter(
		"broker.conversation.started",
		metric.WithDescription("Total number of conversations registered"),
	)
	if err != nil {
		conversationsStarted, _ = noopMeter.Int64Counter("broker.conversation.started")
	}

	conversationsConcluded, err := meter.Int64Counter(
		"broker.conversation.concluded",
		metric.WithDescription("Total number of conversations concluded, by reason"),
	)
	if err != nil {
		conversationsConcluded, _ = noopMeter.Int64Counter("broker.conversation.concluded")
	}

	messagesIngested, err := meter.Int64Counter(
		"broker.message.ingested",
		metric.WithDescription("Total number of messages ingested, by participant role"),
	)
	if err != nil {
		messagesIngested, _ = noopMeter.Int64Counter("broker.message.ingested")
	}

	turnHistogram, err := meter.Float64Histogram(
		"broker.turn.duration",
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(common_metrics.LongLatencySecondsBuckets...),
		metric.WithDescription("Duration of one agent turn, by outcome"),
	)
	if err != nil {
		turnHistogram, _ = noopMeter.Float64Histogram("broker.turn.duration")
	}

	activeTurns, err := meter.Int64UpDownCounter(
		"broker.turn.active",
		metric.WithDescription("Agent turns currently in flight"),
	)
	if err != nil {
		activeTurns, _ = noopMeter.Int64UpDownCounter("broker.turn.active")
	}

	daemonTicks, err := meter.Int64Counter(
		"broker.daemon.tick",
		metric.WithDescription("Processor daemon ticks, by outcome"),
	)
	if err != nil {
		daemonTicks, _ = noopMeter.Int64Counter("broker.daemon.tick")
	}

	agentResolutions, err := meter.Int64Counter(
		"broker.agent.resolution",
		metric.WithDescription("Agent lookups by resolver and outcome"),
	)
	if err != nil {
		agentResolutions, _ = noopMeter.Int64Counter("broker.agent.resolution")
	}

	cacheResults, err := meter.Int64Counter(
		"cache.result",
		metric.WithDescription("Cache lookups by hit/miss"),
	)
	if err != nil {
		cacheResults, _ = noopMeter.Int64Counter("cache.result")
	}

	return &BrokerMetrics{
		conversationsStarted:   conversationsStarted,
		conversationsConcluded: conversationsConcluded,
		messagesIngested:       messagesIngested,
		turnHistogram:          turnHistogram,
		activeTurns:            activeTurns,
		daemonTicks:            daemonTicks,
		agentResolutions:       agentResolutions,
		cacheResults:           cacheResults,
	}
}

func (m *BrokerMetrics) RecordConversationStarted(ctx context.Context) {
	m.conversationsStarted.Add(ctx, 1)
}

func (m *BrokerMetrics) RecordConversationConcluded(ctx context.Context, reason string) {
	m.conversationsConcluded.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func (m *BrokerMetrics) RecordMessageIngested(ctx context.Context, participantRole string) {
	m.messagesIngested.Add(ctx, 1, metric.WithAttributes(attribute.String("participant.role", participantRole)))
}

func (m *BrokerMetrics) RecordTurn(ctx context.Context, outcome string, duration time.Duration) {
	m.turnHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attribute.String("outcome", outcome)))
}

func (m *BrokerMetrics) RecordTurnStarted(ctx context.Context) {
	m.activeTurns.Add(ctx, 1)
}

func (m *BrokerMetrics) RecordTurnFinished(ctx context.Context) {
	m.activeTurns.Add(ctx, -1)
}

func (m *BrokerMetrics) RecordDaemonTick(ctx context.Context, outcome string) {
	m.daemonTicks.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", outcome)))
}

func (m *BrokerMetrics) RecordAgentResolution(ctx context.Context, resolver string, err error) {
	attrs := []attribute.KeyValue{attribute.String("resolver", resolver)}
	// error.type is conditionally required, and only on a failed operation.
	if err != nil {
		attrs = append(attrs, attribute.String("error.type", errorTypeOther))
	}
	m.agentResolutions.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *BrokerMetrics) RecordCacheResult(ctx context.Context, cacheName string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	m.cacheResults.Add(ctx, 1, metric.WithAttributes(
		attribute.String("cache.name", cacheName),
		attribute.String("result", result),
	))
}
